package models

type NewsResponse struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

type Article struct {
	Source struct {
		Name string `json:"name"`
	} `json:"source"`
	Title       string `json:"title"`
	Description string `json:"description"`

}
