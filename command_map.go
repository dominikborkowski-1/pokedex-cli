package main

import "fmt"

// map: pokazuje następne 20 lokacji
func commandMap(cfg *config, args ...string) error {
	resp, err := fetchLocationAreas(cfg.nextLocationsURL)
	if err != nil {
		return err
	}

	cfg.nextLocationsURL = resp.Next
	cfg.prevLocationsURL = resp.Previous

	for _, loc := range resp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

// mapb: pokazuje poprzednie 20 lokacji
func commandMapb(cfg *config, args ...string) error {
	if cfg.prevLocationsURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	resp, err := fetchLocationAreas(cfg.prevLocationsURL)
	if err != nil {
		return err
	}

	cfg.nextLocationsURL = resp.Next
	cfg.prevLocationsURL = resp.Previous

	for _, loc := range resp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}
