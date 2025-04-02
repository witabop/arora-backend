package body

import (
	"arora-search-finger/layer"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
}

func Finger(maxID int64, searchCriteria layer.SearchCriteria) IFinger {
	finger := finger{
		maxID:          maxID,
		validUniverses: &[]layer.Universe{},
		searchCriteria: searchCriteria,
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
		if finger.validateUniverse(universe) {
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

func (finger *finger) validateUniverse(universe layer.Universe) bool {
	critValue := reflect.ValueOf(finger.searchCriteria)
	uniValue := reflect.ValueOf(universe)

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
			return false
		}

		critVal := critField.Elem()
		uniVal := uniField.Elem()

		switch critVal.Kind() {
		case reflect.String:
			if !strings.Contains(uniVal.String(), critVal.String()) {
				return false
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if (uint(float64(critVal.Uint()) * .5)) > uint(uniVal.Uint()) {
				return false
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if (int(float64(critVal.Int()) * .5)) > int(uniVal.Int()) {
				return false
			}

		default:
			if !reflect.DeepEqual(critVal.Interface(), uniVal.Interface()) {
				return false
			}
		}
	}
	return true
}

func (finger *finger) addValidUniverse(vuniverse layer.Universe) {
	finger.mu.Lock()
	defer finger.mu.Unlock()
	*finger.validUniverses = append(*finger.validUniverses, vuniverse)
}
