package models

import "time"

type Driver struct {
	ID                 string       `json:"id"`
	Classification     string       `json:"classification"`
	FirstName          string       `json:"firstName"`
	LastName           string       `json:"lastName"`
	DateOfBirth        time.Time    `json:"dateOfBirth"`
	Email              string       `json:"email"`
	Phone              string       `json:"phone"`
	LicenceNumber      string       `json:"licenceNumber"`
	LicenceIssueDate   time.Time    `json:"licenceIssueDate"`
	LicenceExpiryDate  time.Time    `json:"licenceExpiryDate"`
	LicenceValidUntil  time.Time    `json:"licenceValidUntil"`
	Convictions        []Conviction `json:"convictions"`
}

type Conviction struct {
	ID           string    `json:"id"`
	Date         time.Time `json:"date"`
	OffenceCode  string    `json:"offenceCode"`
	Description  string    `json:"description"`
	PenaltyPoints int      `json:"penaltyPoints"`
	FineAmount   float64   `json:"fineAmount"`
}
