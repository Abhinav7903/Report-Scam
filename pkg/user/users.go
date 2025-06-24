package users

import "abuse/factory"

type Repository interface {
	//Create User
	CreateUser(data factory.User) (string, error)
	//Get User
	GetUser(email string) (factory.User, error)
	// is Reporter or User
	IsReporter(email string) (bool, error)
	IncrementReporterReportCount(reporterID string) error
	// VerifyEmail verifies the user's email
	VerifyEmail(email string) error
}
