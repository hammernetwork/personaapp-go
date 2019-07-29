package models

// Role describes a role of the user
type Role string

const (
	// CompanyRole is company role
	CompanyRole Role = "company"
	// EmployeeRole is employee role
	EmployeeRole Role = "employee"
)
