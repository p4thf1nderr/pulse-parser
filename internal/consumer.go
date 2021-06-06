package internal

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/p4thf1nderr/pulse-parser/model"
	"golang.org/x/net/html"
)

type Consumer struct{}

func NewConsumer() *Consumer {
	return &Consumer{}
}

type row struct {
	Url        string   `json:"url"`
	Categories []string `json:"categories"`
}

func (c *Consumer) Consume(orgn chan []string, crBuffer chan map[string][]model.Record) {
	// run consumer

	for {
		select {
		case messages, ok := <-orgn:
			if ok {
				cat := make(map[string][]model.Record)
				for _, item := range messages {
					var m row
					json.Unmarshal([]byte(item), &m)
					for _, c := range m.Categories {
						cat[c] = append(cat[c], model.Record{
							Url: m.Url,
						})
					}
				}
				crBuffer <- cat
			} else {
				close(crBuffer)
				return
			}
		}
	}
}

func (c *Consumer) Fetch(ch <-chan map[string][]model.Record, stop chan struct{}) {
	defer func() {
		close(stop)
	}()

	var client = &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	var counter int

	for {
		select {
		case val, ok := <-ch:
			if ok {
				for k, v := range val {
					writerBuffer := make(chan *model.Record)

					for _, rec := range v {
						go func(rec model.Record) {
							hm := new(model.Record)
							hm.Url = rec.Url

							resp, err := client.Get(rec.Url)
							if err != nil {
								fmt.Printf("error occured: %v\n", err)
								writerBuffer <- hm
							} else {
								parse(resp, hm)
								writerBuffer <- hm
							}
						}(rec)
					}

					for range v {
						hm := <-writerBuffer

						if hm != nil && hm.Url != "" && hm.Description != "" && hm.Title != "" {
							counter++
							f, err := os.OpenFile(fmt.Sprintf("data/output/%s.txt", k), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
							if err != nil {
								fmt.Printf("error occured:%v\n", err)
							}
							f.WriteString(fmt.Sprintf("%s\\%s\\%s\n", hm.Url, hm.Title, hm.Description))
							fmt.Printf("\r %d parsed", counter)
						}
					}
				}
			} else {
				return
			}
		}
	}
}

func extractMetaProperty(t html.Token, prop string) (content string, ok bool) {
	for _, attr := range t.Attr {
		if attr.Key == "property" && attr.Val == prop {
			ok = true
		}

		if attr.Key == "content" {
			content = attr.Val
		}
	}
	return
}

func parse(response *http.Response, hm *model.Record) {

	z := html.NewTokenizer(response.Body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := z.Token()

			if t.Data == "meta" {
				ogTitle, ok := extractMetaProperty(t, "og:title")
				if ok {
					hm.Title = ogTitle
				}

				ogDesc, ok := extractMetaProperty(t, "og:description")
				if ok {
					hm.Description = ogDesc
				}
			}
		}
	}
}
