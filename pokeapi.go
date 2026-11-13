package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"pokedexcli/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2"

// cacheInterval: jak długo wpis żyje w cache. 5 minut wystarczy,
// żeby swobodnie skakać map/mapb bez ponownych zapytań.
const cacheInterval = 5 * time.Minute

// cache tworzony raz przy starcie programu i współdzielony przez wszystkie
// zapytania. Kluczem jest pełny URL, wartością surowe bajty odpowiedzi.
var cache = pokecache.NewCache(cacheInterval)

// Odpowiedź z endpointu /location-area.
// Next i Previous to *string, bo API zwraca null na pierwszej/ostatniej stronie.
type locationAreasResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// Odpowiedź z endpointu /location-area/{name}.
// API zwraca dużo więcej pól, ale json.Unmarshal ignoruje te,
// których nie ma w structcie, więc bierzemy tylko to, czego potrzebujemy.
type locationAreaResponse struct {
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

// fetchJSON pobiera url i dekoduje JSON do `out` (musi być wskaźnikiem).
// Najpierw sprawdza cache; przy cache miss robi zapytanie HTTP
// i zapisuje bajty w cache. Wspólne dla wszystkich endpointów.
func fetchJSON(url string, out any) error {
	// Cache hit: dekodujemy zapisane bajty, bez zapytania HTTP.
	if data, ok := cache.Get(url); ok {
		return json.Unmarshal(data, out)
	}

	// Cache miss: zwykłe zapytanie.
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found")
	}
	if res.StatusCode > 299 {
		return fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, out); err != nil {
		return err
	}

	// Zapisujemy dopiero gdy status był OK i JSON się zdekodował,
	// żeby nie cache'ować błędów.
	cache.Add(url, data)
	return nil
}

// fetchLocationAreas pobiera jedną stronę location areas.
// Jeśli pageURL jest nil, pobiera pierwszą stronę.
func fetchLocationAreas(pageURL *string) (locationAreasResponse, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	var resp locationAreasResponse
	err := fetchJSON(url, &resp)
	return resp, err
}

// fetchLocationArea pobiera szczegóły jednej lokacji (m.in. listę Pokemonów).
func fetchLocationArea(name string) (locationAreaResponse, error) {
	url := baseURL + "/location-area/" + name

	var resp locationAreaResponse
	err := fetchJSON(url, &resp)
	return resp, err
}

// Odpowiedź z endpointu /pokemon/{name}. base_experience służy do szansy
// złapania; reszta pól jest wyświetlana przez `inspect`.
type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	// W JSON-ie każda statystyka to {"base_stat": 45, "stat": {"name": "attack"}}.
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	// Analogicznie typy: {"slot": 1, "type": {"name": "normal"}}.
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

// fetchPokemon pobiera dane jednego Pokemona po nazwie.
func fetchPokemon(name string) (Pokemon, error) {
	url := baseURL + "/pokemon/" + name

	var resp Pokemon
	err := fetchJSON(url, &resp)
	return resp, err
}
