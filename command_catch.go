package main

import (
	"errors"
	"fmt"
	"math/rand"
)

// catchThreshold: rzut to losowa liczba z [0, base_experience).
// Złapany, jeśli wypadnie poniżej progu. Im wyższe base_experience,
// tym mniejsza szansa: pikachu (112) ~36%, mewtwo (340) ~12%.
const catchThreshold = 40

// catch <pokemon>: rzuca Pokeballem; przy sukcesie dodaje Pokemona do pokedexu.
func commandCatch(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("usage: catch <pokemon>")
	}
	name := args[0]

	// Wymaganie z lekcji: ten komunikat przed ustaleniem wyniku.
	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	pokemon, err := fetchPokemon(name)
	if err != nil {
		return err
	}

	// rand.Intn(n) panikuje dla n <= 0, więc zabezpieczamy brak danych.
	if pokemon.BaseExperience > 0 && rand.Intn(pokemon.BaseExperience) >= catchThreshold {
		fmt.Printf("%s escaped!\n", name)
		return nil
	}

	fmt.Printf("%s was caught!\n", name)
	cfg.pokedex[name] = pokemon
	return nil
}
