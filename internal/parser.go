package internal

import (
	"github.com/p4thf1nderr/pulse-parser/model"
	"github.com/p4thf1nderr/pulse-parser/usecase"
)

type Parser struct {
	Producer usecase.Producer
	Consumer usecase.Consumer

	readBuffer     chan string
	consumerBuffer chan []string
	crawlerBuffer  chan map[string][]model.Record
}

func NewParser(producer usecase.Producer, consumer usecase.Consumer) *Parser {
	return &Parser{
		Producer:       producer,
		Consumer:       consumer,
		readBuffer:     make(chan string, 10),
		consumerBuffer: make(chan []string),
		crawlerBuffer:  make(chan map[string][]model.Record),
	}
}

func (p *Parser) Run() {
	file := "data/500.jsonl"

	stop := make(chan struct{})

	go p.Producer.Produce(p.consumerBuffer, p.readBuffer)
	go p.Consumer.Fetch(p.crawlerBuffer, stop)
	go p.Consumer.Consume(p.consumerBuffer, p.crawlerBuffer)

	p.Producer.ReadInput(file, p.readBuffer)

	for {
		select {
		case <-stop:
			return
		}
	}
	//close(stop)
}
