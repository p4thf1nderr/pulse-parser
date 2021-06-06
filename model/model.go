package model

type Record struct {
	Url         string
	Title       string `json:"title"`
	Description string `json:"description"`
}
