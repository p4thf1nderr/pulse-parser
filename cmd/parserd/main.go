package main

import (
	"github.com/p4thf1nderr/pulse-parser/internal"
)

type HTMLMeta struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {
	// читаем из файла берем урл и массив категорий
	// на каждую категорию должен быть отдельный файл в выводе

	// последовательно обрабатываем адреса из файла
	// формируем буфер где будут хранится данные вида
	// название категории - структура {массив ссылок}

	producer := internal.NewProducer()
	consumer := internal.NewConsumer()

	parser := internal.NewParser(producer, consumer)
	parser.Run()
}