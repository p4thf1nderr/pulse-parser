package usecase

import "github.com/p4thf1nderr/pulse-parser/model"

type Consumer interface {
	Consume(chan []string, chan map[string][]model.Record)
	Fetch(<-chan map[string][]model.Record, chan struct{})
}

type Producer interface {
	Produce(chan []string, chan string)
	ReadInput(string, chan string)
}
