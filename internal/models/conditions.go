package models

type Conditions struct {
	ID       string  `json:"id"`
	Temp     float64 `json:"temp"`
	Humidity float64 `json:"humidity"`
	Time     string  `json:"time"`
}
