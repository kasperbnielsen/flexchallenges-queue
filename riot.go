package main

import (
	"fmt"
	"io"
	"net/http"
)

var APIKEY string = "RGAPI-e665e95c-a5c5-4b4e-bcfc-9e48c8703024"

func getMatch(matchid string) []byte {
	url := fmt.Sprintf("https://europe.api.riotgames.com/lol/match/v5/matches/%s?api_key=%s", matchid, APIKEY)

	response, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	return body
}

func getMatchIds(puuid string, time string) []byte {
	url := fmt.Sprintf("https://europe.api.riotgames.com/lol/match/v5/matches/by-puuid/%s/ids?startTime=%s&type=ranked&count=100&api_key=%s", puuid, time, APIKEY)

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

	return body
}
