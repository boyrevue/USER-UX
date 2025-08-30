package models

import "time"

type Session struct {
	ID        string    `json:"id"`
	Language  string    `json:"language"`
	Drivers   []Driver  `json:"drivers"`
	Vehicles  []Vehicle `json:"vehicles"`
	Claims    Claims    `json:"claims"`
	Policy    Policy    `json:"policy"`
	Documents []Document `json:"documents"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Claims struct {
	Claims    []Claim    `json:"claims"`
	Accidents []Accident `json:"accidents"`
}

type Claim struct {
	ID          string  `json:"id"`
	Date        string  `json:"date"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Settled     bool    `json:"settled"`
}

type Accident struct {
	ID            string  `json:"id"`
	Date          string  `json:"date"`
	Type          string  `json:"type"`
	Description   string  `json:"description"`
	EstimatedCost float64 `json:"estimatedCost"`
	FaultClaim    bool    `json:"faultClaim"`
}

type Policy struct {
	StartDate    string  `json:"startDate"`
	CoverType    string  `json:"coverType"`
	Excess       float64 `json:"excess"`
	NCDYears     int     `json:"ncdYears"`
	NCDProtected bool    `json:"ncdProtected"`
}

type Document struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Path string `json:"path"`
}
