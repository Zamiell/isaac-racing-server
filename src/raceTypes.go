package server

type RacerAddItemMessage struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Item *Item  `json:"item"`
}
