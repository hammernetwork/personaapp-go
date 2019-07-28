package models

// Role describes a role of the user
type Role string

const (
	// Company is company role
	Company Role = "company"
	// Employee is employee role
	Employee Role = "employee"
)
