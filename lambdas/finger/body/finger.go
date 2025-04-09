package body

import (
	"arora-search-finger/layer"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

const (
	BASE_URL     string = "https://games.roblox.com/v1/games?universeIds="
	MAX_APIS     int8   = 10 // Max known APIs possible from a single IP
	ID_INCREMENT int64  = 49 // How many ids to search per api call (add one because of start inclusion for number of actual ids)
)

var client http.Client

type IFinger interface {
	Feel() []layer.Universe
}

type finger struct {
	maxID          int64
	validUniverses *[]layer.Universe
	searchCriteria layer.SearchCriteria
	mu             sync.Mutex
	prod           bool
}

func Finger(maxID int64, searchCriteria layer.SearchCriteria, stage string) IFinger {
	finger := finger{
		maxID:          maxID,
		validUniverses: &[]layer.Universe{},
		searchCriteria: searchCriteria,
		prod:           stage == "prod",
	}
	return &finger
}

// func (finger *finger) Feel() []int64 {
func (finger *finger) Feel() []layer.Universe {
	var wg sync.WaitGroup

	var i int8
	start := rand.Int63n(finger.maxID) + 1
	for i = 0; i < MAX_APIS; i++ {
		wg.Add(1)
		go func(s int64) {
			defer wg.Done()
			finger.touch(s, s+ID_INCREMENT)
		}(start)
		start += 50
	}
	wg.Wait()

	return *finger.validUniverses
}

func (finger *finger) touch(start int64, end int64) {
	log.Printf("Searching IDs %v-%v\n", start, end)
	defer log.Printf("Finished Searching IDs %v-%v\n", start, end)

	resp, err := finger.touchIDs(start, end)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var universeResponse layer.UniverseResponse
	if err := json.Unmarshal(body, &universeResponse); err != nil {
		return
	}

	if universeResponse.Data == nil {
		return
	}

	for _, universe := range *universeResponse.Data {
		finger.calculatePercentage(&universe)
		if !finger.prod || *universe.PercentMatch >= 0.5 {
			finger.addValidUniverse(universe)
		}
	}
}

func (finger *finger) touchIDs(start int64, end int64) (*http.Response, error) {
	url := fmt.Sprintf("%v%v", BASE_URL, finger.formatIDs(start, end))
	req, err := http.NewRequest("GET", url, &bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (finger *finger) formatIDs(start int64, end int64) string {
	var builder strings.Builder

	for i := start; i <= end; i++ {
		builder.WriteString(strconv.FormatInt(i, 10))
		if i < end {
			builder.WriteString(",")
		}
	}

	return builder.String()
}

func (finger *finger) calculatePercentage(universe *layer.Universe) {
	critValue := reflect.ValueOf(finger.searchCriteria)
	uniValue := reflect.ValueOf(universe).Elem()
	percentages := []float64{}

	for i := 0; i < critValue.NumField(); i++ {
		critField := critValue.Field(i)
		if critField.IsNil() {
			continue
		}

		if uniValue.Kind() == reflect.Ptr {
			uniValue = uniValue.Elem()
		}

		fieldName := critValue.Type().Field(i).Name

		uniField := uniValue.FieldByName(fieldName)
		if !uniField.IsValid() {
			continue
		}

		critVal := critField.Elem()
		uniVal := uniField.Elem()

		percentages = append(percentages, CompareValues(critVal, uniVal))
	}

	if universe.PercentMatch == nil {
		universe.PercentMatch = new(float64)
	}
	*universe.PercentMatch = math.Round(Avg(percentages)*10000) / 10000
}

func (finger *finger) addValidUniverse(vuniverse layer.Universe) {
	finger.mu.Lock()
	defer finger.mu.Unlock()
	*finger.validUniverses = append(*finger.validUniverses, vuniverse)
}
