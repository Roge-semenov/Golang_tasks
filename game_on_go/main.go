package main

type Player struct {
	Inventory             map[string]bool
	WearingBackpack       bool
	CurrentPlace          *Place
	externalitemjustforPR UiT
}
