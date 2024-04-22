package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type cliCommand struct {
	name string
	description string 
	callback func() error
}

type Location struct {
	Count    int    `json:"count"`
	Next     *string `json:"next"`
	Previous *string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type mapHistory struct {
	Base string
	Next *string
	Previous *string
}



func main() {
	startRepl()
}

func startRepl() {
	history := &mapHistory{Base:"https://pokeapi.co/api/v2/location", Previous: nil, Next: nil}
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("pokedex > ")
		scanner.Scan()
		commandName := scanner.Text()
		if len(commandName) == 0 {
			continue
		}
		commandMap := getCommands(history)
		command, ok := commandMap[commandName]
		if !ok {
			fmt.Println("Invalid command")
			continue
		}
		command.callback()
	} 
}

func mapCommand(history *mapHistory) error {
	url := ""
	if history.Next == nil {
		url = history.Base
	} else {
		url = *history.Next
	}

	location, err := getRequest(url)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	history.Previous = location.Previous
	history.Next = location.Next
	
	for _, data := range location.Results {
		fmt.Println(data.Name)
	}

	return nil 
	
}

func mapbCommand(history *mapHistory) error {
	url := ""
	if history.Previous == nil {
		url = history.Base
	} else {
		url = *history.Previous
	}

	location, err := getRequest(url)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	if history.Previous != nil {
		history.Previous = location.Previous
	}
	history.Next = location.Next
	
	for _, data := range location.Results {
		fmt.Println(data.Name)
	}

	return nil 
	
}


func getRequest(url string)  (Location, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
		return Location{}, err
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 {
		log.Fatalf("Response failed with status code %d and \n body :%s\n", resp.StatusCode, body)
		return Location{}, err
	}
	if err != nil {
		log.Fatal(err)
		return Location{}, err
	}
	location := Location{}
	err = json.Unmarshal(body, &location)
	if err != nil {
		return location, err
	}
	return location, nil
}

func getCommands(history *mapHistory) map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    help,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    exit,
		},
		"map": {
			name: "map", 
			description: "Get the next page of locations",
			callback: func() error {
				return mapCommand(history)
			},
		},

		"mapb": {
			name: "mapb", 
			description: "Get the previous page of locations",
			callback: func() error {
				return mapbCommand(history)
			},
		},
		
	}
}



func help() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: ")
	commands := getCommands(nil)
	for _, cmd := range commands{
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}


func exit() error {
	os.Exit(0)
	return nil
}