package body

import (
	"arora-search-brain/layer"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	MAX_ACTIVE_NEURONS uint8  = 50
	FINGER_URL         string = "https://r8oxhje7na.execute-api.us-east-1.amazonaws.com/dev/search/finger"
)

var client http.Client

type IBrain interface {
	Think() []layer.Universe
}

type brain struct {
	validUniverses []layer.Universe
	numGames       uint8
	searchCriteria layer.SearchCriteria
	mu             sync.Mutex
	maxID          int64
}

func Brain(numGames uint8, searchCriteria layer.SearchCriteria) IBrain {
	brain := brain{
		validUniverses: []layer.Universe{},
		numGames:       numGames,
		searchCriteria: searchCriteria,
		maxID:          7_000_000_000,
	}
	return &brain
}

func (brain *brain) Think() []layer.Universe {
	log.Println("Thinking for universes...")
	defer log.Printf("Finished Thinking!!!")
	var (
		neuronCtr uint8 = 0
		neuronCh        = make(chan []layer.Universe, MAX_ACTIVE_NEURONS)
		timeout         = time.After(28 * time.Second)
	)

	for len(brain.validUniverses) < int(brain.numGames) {
		select {
		case <-timeout:
			log.Println("Could not get numGames in alloted time")
			return brain.validUniverses

		default:
			for neuronCtr < MAX_ACTIVE_NEURONS {
				go brain.prepareNeuron(neuronCh)
				neuronCtr++
			}
			universes := <-neuronCh
			if len(universes) > 0 {
				brain.addValidUniverses(universes)
			}
			go brain.prepareNeuron(neuronCh)
		}
	}
	return brain.validUniverses
}

func (brain *brain) prepareNeuron(resultCh chan<- []layer.Universe) {
	universes := []layer.Universe{}

	resp, err := brain.activateNeuron()
	if err != nil {
		resultCh <- universes
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		resultCh <- universes
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		resultCh <- universes
		return
	}

	var universeResponse layer.UniverseResponse
	if err := json.Unmarshal(body, &universeResponse); err != nil {
		resultCh <- universes
		return
	}

	if universeResponse.ValidUniverses == nil {
		resultCh <- universes
		return
	}

	universes = *universeResponse.ValidUniverses
	resultCh <- universes
}

func (brain *brain) activateNeuron() (*http.Response, error) {
	payload := map[string]interface{}{
		"maxID":          brain.maxID,
		"searchCriteria": brain.searchCriteria,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", FINGER_URL, bytes.NewBuffer(jsonData))
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

func (brain *brain) addValidUniverses(universes []layer.Universe) {
	brain.mu.Lock()
	defer brain.mu.Unlock()
	addCount := int(brain.numGames) - len(brain.validUniverses)
	if addCount <= 0 {
		return
	}
	if len(universes) < addCount {
		addCount = len(universes)
	}
	brain.validUniverses = append(brain.validUniverses, universes[:addCount]...)
}
