package models

type Employee struct {
	FullName         string `json:"full_name"`
	Preferred        string `json:"preferred"`
	Email            string `json:"email"`
	UniqueIdentifier string `json:"unique_identifier"`
	ManagersEmail    string `json:"managers_email"`
	StartDate        string `json:"start_date"`
	Tenure           string `json:"tenure"`
	Language         string `json:"language"`
}
