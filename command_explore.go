package main

import (
	"errors"
	"fmt"
)

// explore <area>: pokazuje Pokemony, które można spotkać w danej lokacji.
// Nazwa lokacji przychodzi jako pierwszy argument (np. z wyniku `map`).
func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("usage: explore <location-area>")
	}
	name := args[0]

	fmt.Printf("Exploring %s...\n", name)
	area, err := fetchLocationArea(name)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, enc := range area.PokemonEncounters {
		fmt.Printf(" - %s\n", enc.Pokemon.Name)
	}
	return nil
}
