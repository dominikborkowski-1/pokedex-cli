package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type config struct {
	commandRegistry  map[string]cliCommand
	nextLocationsURL *string
	prevLocationsURL *string
	// Złapane Pokemony, klucz to nazwa.
	pokedex map[string]Pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func getCommands() map[string]cliCommand{
	return map[string]cliCommand{
    "exit": {
        name:        "exit",
        description: "Exit the Pokedex",
        callback:    commandExit,
    },
	"help": {
		name: "help",
		description: "Displays a help message",
		callback: help,
	},
	"map": {
		name:        "map",
		description: "Displays the next 20 location areas",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Displays the previous 20 location areas",
		callback:    commandMapb,
	},
	"explore": {
		name:        "explore",
		description: "Lists Pokemon in a location area: explore <area-name>",
		callback:    commandExplore,
	},
	"catch": {
		name:        "catch",
		description: "Try to catch a Pokemon: catch <pokemon-name>",
		callback:    commandCatch,
	},
	"inspect": {
		name:        "inspect",
		description: "Show details of a caught Pokemon: inspect <pokemon-name>",
		callback:    commandInspect,
	},
}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}
		cmd, ok := cfg.commandRegistry[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		if err := cmd.callback(cfg, words[1:]...); err != nil {
			fmt.Println(err)
		}
	}
}