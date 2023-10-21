package main

import (
	"net/http"
	"fmt"
	"io"
	"encoding/json"
)

func getMatchIds(puuid string, time string) string {
	APIKEY := "RGAPI-e665e95c-a5c5-4b4e-bcfc-9e48c8703024"
	url := "https://europe.api.riotgames.com/lol/match/v5/matches/by-puuid/" + string(puuid) + "/ids?startTime=" + time + "&type=ranked&count=100&api_key=" + APIKEY

	fmt.Println(url)

	response, err := http.Get(url)

	if err != nil {
		panic(err)
	}


	
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	var temp []string

	newbody := json.Unmarshal([]byte(response.Body), &temp)

	fmt.Println(newbody)

	return string(body)
}