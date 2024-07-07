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

func initGame() {
	kitchen = Place{
		Description: "ты находишься на кухне, на столе: чай, надо собрать рюкзак и идти в универ. можно пройти - коридор",
		Items:       map[string]string{"чай": "на столе"},
		Exits:       map[string]*Place{"коридор": &corridor},
	}

	corridor = Place{
		Description: "ничего интересного. можно пройти - кухня, комната, улица",
		Exits:       map[string]*Place{"кухня": &kitchen, "комната": &room, "улица": &street},
	}

	room = Place{
		Description: "ты в своей комнате. можно пройти - коридор",
		Items:       map[string]string{"ключи": "на столе", "конспекты": "на столе", "рюкзак": "на стуле"},
		Exits:       map[string]*Place{"коридор": &corridor},
	}

	street = Place{
		Description: "на улице весна. можно пройти - домой",
	}

	player = Player{
		Inventory:    make(map[string]bool),
		CurrentPlace: &kitchen,
	}
}

func lookAround(player *Player) string {
	description := player.CurrentPlace.Description

	var itemsDescription []string

	if len(player.CurrentPlace.Items) > 0 {
		var onTable, onChair []string
		for item, location := range player.CurrentPlace.Items {
			if location == "на столе" {
				onTable = append(onTable, item)
			} else if location == "на стуле" {
				onChair = append(onChair, item)
			}
		}
		if len(onTable) > 0 {
			itemsDescription = append(itemsDescription, "на столе: "+strings.Join(onTable, ", "))
		}
		if len(onChair) > 0 {
			itemsDescription = append(itemsDescription, "на стуле: "+strings.Join(onChair, ", "))
		}
	}

	if len(itemsDescription) > 0 {
		description = strings.Join(itemsDescription, ". ") + ". " + description
	}

	return description
}

func handleCommand(command string) string {
	parts := strings.Split(command, " ")
	if len(parts) == 0 {
		return "неизвестная команда"
	}

	switch parts[0] {
	case "осмотреться":
		return lookAround(&player)
	case "идти":
		if len(parts) < 2 {
			return "куда идти?"
		}
		direction := parts[1]
		nextPlace, ok := player.CurrentPlace.Exits[direction]
		if !ok {
			return fmt.Sprintf("нет пути в %s", direction)
		}
		if player.CurrentPlace == &corridor && direction == "улица" && !player.Inventory["ключи"] {
			return "дверь закрыта"
		}
		player.CurrentPlace = nextPlace
		if player.CurrentPlace == &kitchen && direction == "кухня" {
			return "кухня, ничего интересного. можно пройти - коридор"
		}
		return nextPlace.Description
	case "надеть":
		if len(parts) < 2 {
			return "что надеть?"
		}
		item := parts[1]
		if item == "рюкзак" {
			player.WearingBackpack = true
			return "вы надели: рюкзак"
		}
		return "неизвестный предмет"
	case "взять":
		if len(parts) < 2 {
			return "что взять?"
		}
		item := parts[1]
		if player.WearingBackpack {
			if _, ok := player.CurrentPlace.Items[item]; ok {
				player.Inventory[item] = true
				delete(player.CurrentPlace.Items, item)
				return fmt.Sprintf("предмет добавлен в инвентарь: %s", item)
			}
			return "нет такого"
		}
		return "некуда класть"
	case "применить":
		if len(parts) < 3 {
			return "не к чему применить"
		}
		item := parts[1]
		target := parts[2]
		if player.Inventory[item] {
			if target == "дверь" && item == "ключи" {
				return "дверь открыта"
			}
			return "не к чему применить"
		}
		return fmt.Sprintf("нет предмета в инвентаре - %s", item)
	default:
		return "неизвестная команда"
	}
}

func (r *Place) AddAction(command string, action func(*Player) string) {
	if r.Actions == nil {
		r.Actions = make(map[string]func(*Player) string)
	}
	r.Actions[command] = action
}
