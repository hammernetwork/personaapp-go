package validator

type Company struct {
	Name string
	Email string
	Phone string
	Password string
}

func ValidateCompany(c Company) error {
	// validate name
	// validate email
	// validate phone
	// validate password
	return nil
}

func validateName(name string) error {
	// validate length
	// validate allowed symbols
	return nil
}

func validateEmail(name string) error {
	// validate length
	// validate format
	return nil
}

func validatePhone(name string) error {
	// validate length
	// validate format
	return nil
}

func validatePassword(name string) error {
	// validate length
	// validate allowed symbols
	return nil
}