package main

type Observer interface {
	NotifyObserver(field string, data string)
}

type Observable interface {
	RegisterObserver(obs Observer)
}
