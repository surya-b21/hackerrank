package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	start := time.Now()
	totalPage := int64(0)

	response, _ := http.Get("https://jsonmock.hackerrank.com/api/football_matches?year=2011")

	responseData, _ := ioutil.ReadAll(response.Body)

	data, _ := UnmarshalResponseModel(responseData)
	totalPage = data.TotalPages

	totalDrawn := make(chan int64, totalPage)
	for i := 1; i <= int(totalPage); i++ {
		wg.Add(1)
		go getMatchPerPage(i, totalDrawn, &wg)
	}

	go func() {
		wg.Wait()
		close(totalDrawn)
	}()

	totalDrawnInt := int64(0)
	for drawn := range totalDrawn {
		totalDrawnInt += drawn
	}
	fmt.Println(totalDrawnInt)

	elapsed := time.Since(start)
	fmt.Printf("Process took %s", elapsed)
}

func getMatchPerPage(page int, channel chan int64, wg *sync.WaitGroup) {
	defer wg.Done()

	drawnMath := int64(0)
	response, _ := http.Get(fmt.Sprintf("https://jsonmock.hackerrank.com/api/football_matches?year=2011&page=%d", page))

	responseData, _ := ioutil.ReadAll(response.Body)

	data, _ := UnmarshalResponseModel(responseData)

	for _, match := range data.Data {
		if match.Team1Goals == match.Team2Goals {
			drawnMath++
		}
	}

	channel <- drawnMath
}

func UnmarshalResponseModel(data []byte) (ResponseModel, error) {
	var r ResponseModel
	err := json.Unmarshal(data, &r)
	return r, err
}

type ResponseModel struct {
	Page       int64   `json:"page"`
	PerPage    int64   `json:"per_page"`
	Total      int64   `json:"total"`
	TotalPages int64   `json:"total_pages"`
	Data       []Datum `json:"data"`
}

type Datum struct {
	Competition string `json:"competition"`
	Year        int64  `json:"year"`
	Round       string `json:"round"`
	Team1       string `json:"team1"`
	Team2       string `json:"team2"`
	Team1Goals  string `json:"team1goals"`
	Team2Goals  string `json:"team2goals"`
}
