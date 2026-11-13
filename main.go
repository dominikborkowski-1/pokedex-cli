package main

func main() {
	commands := getCommands()
	var cfg config
	cfg.commandRegistry = commands
	cfg.pokedex = make(map[string]Pokemon)
	startRepl(&cfg)
}