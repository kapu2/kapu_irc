package main

type Model struct {
	channel string
}

func NewModel() *Model {
	return &Model{}
}

func (m *Model) SetChannel(channel string) {
	m.channel = channel
}

func (m *Model) GetChannel() string {
	return m.channel
}
