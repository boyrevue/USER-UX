package models

type Vehicle struct {
	ID              string   `json:"id"`
	Registration    string   `json:"registration"`
	Make            string   `json:"make"`
	Model           string   `json:"model"`
	Year            int      `json:"year"`
	EngineSize      string   `json:"engineSize"`
	FuelType        string   `json:"fuelType"`
	Transmission    string   `json:"transmission"`
	EstimatedValue  float64  `json:"estimatedValue"`
	Modifications   []string `json:"modifications"`
}
