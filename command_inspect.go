package main

import (
	"errors"
	"fmt"
)

// inspect <pokemon>: pokazuje szczegóły złapanego Pokemona.
// Bez zapytania do API — dane zostały zapisane w pokedexie przy `catch`.
func commandInspect(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("usage: inspect <pokemon>")
	}
	name := args[0]

	pokemon, ok := cfg.pokedex[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, s := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
	return nil
}
