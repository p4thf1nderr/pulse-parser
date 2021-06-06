package internal

import (
	"bufio"
	"log"
	"os"
)

type Producer struct{}

func NewProducer() *Producer {
	return &Producer{}
}

func (p *Producer) Produce(dst chan []string, orgn chan string) {

	var msgs []string

	defer close(dst)

	for {
		select {
		case message, ok := <-orgn:
			if ok {
				msgs = append(msgs, message)
				if len(msgs) >= 10 {
					dst <- msgs
					msgs = []string{}
				}
			} else {
				return
			}
		}
	}
}

func (p *Producer) ReadInput(path string, dst chan string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dst <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	close(dst)
}
