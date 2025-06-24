package cases

import "abuse/factory"

type Repository interface {
	// CreateCase creates a new case with the given data.
	CreateCase(data factory.Report) (string, error)
	// GetAllCases retrieves all cases.
	GetAllCases() ([]factory.Report, error)
	// GetCasebyUser retrieves cases by user email.
	GetCasebyUser(email string) ([]factory.Report, error)
}
