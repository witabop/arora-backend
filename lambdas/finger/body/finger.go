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
var suc int8 = 0

type IFinger interface {
	// Feel() []int64
	Feel() ([]int64, int8)
}

type finger struct {
	maxID    int64
	validIDs *[]int64
	mu       sync.Mutex
}

func Finger(maxID int64) IFinger {
	finger := finger{
		maxID:    maxID,
		validIDs: &[]int64{},
	}
	return &finger
}

// func (finger *finger) Feel() []int64 {
func (finger *finger) Feel() ([]int64, int8) {
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

	log.Println(suc)

	return *finger.validIDs, suc
}

func (finger *finger) touch(start int64, end int64) {
	log.Printf("Searching IDs %v-%v\n", start, end)
	defer log.Printf("Finished Searching IDs %v-%v\n", start, end)

	resp, err := finger.touchIDs(start, end)
	if err != nil || resp.StatusCode != http.StatusOK {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var UniverseResponse layer.UniverseResponse
	if err := json.Unmarshal(body, &UniverseResponse); err != nil {
		return
	}

	if UniverseResponse.Data == nil {
		return
	}

	finger.mu.Lock()
	suc++
	finger.mu.Unlock()

	for _, universe := range *UniverseResponse.Data {
		if finger.validateUniverse(universe) {
			finger.addValidID(*universe.RootPlaceID)
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
	return (*universe.Playing >= 1 || *universe.Visits >= 1)
}

func (finger *finger) addValidID(vid int64) {
	finger.mu.Lock()
	defer finger.mu.Unlock()
	*finger.validIDs = append(*finger.validIDs, vid)
}
