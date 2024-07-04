package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Player struct {
	Inventory       map[string]bool
	WearingBackpack bool
	CurrentPlace    *Place
}

type Place struct {
	Description string
	Items       map[string]string
	Exits       map[string]*Place
	Actions     map[string]func(*Player) string
}

func main() {
	initGame()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Добро пожаловать в игру! Введите команду:")

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "выход" {
			break
		}

		output := handleCommand(text)
		fmt.Println(output)
	}
}

var player Player
var kitchen, corridor, room, street Place

func (r *Place) AddAction(command string, action func(*Player) string) {
	if r.Actions == nil {
		r.Actions = make(map[string]func(*Player) string)
	}
	r.Actions[command] = action
}
