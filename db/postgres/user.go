package postgres

import (
	"abuse/factory"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (p *Postgres) CreateUser(data factory.User) (string, error) {

	if data.ID == "" {
		data.ID = uuid.New().String()
	}

	switch data.UserType {
	case "user":
		query := `INSERT INTO users (id, email, name, contact_info, created_at) 
		          VALUES ($1, $2, $3, $4, NOW())`
		_, err := p.dbConn.Exec(query, data.ID, data.Email, data.Name, data.ContactInfo)
		if err != nil {
			logrus.Error("Error creating user:", err)
			return "", err
		}
		logrus.Info("User created successfully:", data.Email)
		return data.ID, nil

	case "reporter":
		query := `INSERT INTO reporters (id, email, name, affiliation, created_at) 
		          VALUES ($1, $2, $3, $4, NOW())`
		_, err := p.dbConn.Exec(query, data.ID, data.Email, data.Name, data.Affiliation)
		if err != nil {
			logrus.Error("Error creating reporter:", err)
			return "", err
		}
		logrus.Info("Reporter created successfully:", data.Email)
		return data.ID, nil

	default:
		err := fmt.Errorf("invalid user type: %s", data.UserType)
		logrus.Error(err)
		return "", err
	}
}

func (p *Postgres) VerifyEmail(email string) error {
	// First, try to update users
	res, err := p.dbConn.Exec(`UPDATE users SET isVerified = TRUE WHERE email = $1`, email)
	if err != nil {
		return fmt.Errorf("failed to verify email in users: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected > 0 {
		return nil // Email found and verified in users
	}

	// Try reporters next
	res, err = p.dbConn.Exec(`UPDATE reporters SET isVerified = TRUE WHERE email = $1`, email)
	if err != nil {
		return fmt.Errorf("failed to verify email in reporters: %w", err)
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected for reporters: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("email %s not found in users or reporters", email)
	}

	return nil
}

func (p *Postgres) GetUser(email string) (factory.User, error) {
	var user factory.User

	// Check in users table
	queryUser := `SELECT id, email, name, contact_info, created_at FROM users WHERE email = $1`
	err := p.dbConn.QueryRow(queryUser, email).Scan(&user.ID, &user.Email, &user.Name, &user.ContactInfo, &user.CreatedAt)
	if err == nil {
		user.UserType = "user"
		logrus.Info("User retrieved from users table:", user.Email)
		return user, nil
	}

	// Check in reporters table
	queryReporter := `SELECT id, email, name, affiliation, total_reports, created_at FROM reporters WHERE email = $1`
	err = p.dbConn.QueryRow(queryReporter, email).Scan(&user.ID, &user.Email, &user.Name, &user.Affiliation, &user.TotalReports, &user.CreatedAt)
	if err == nil {
		user.UserType = "reporter"
		logrus.Info("User retrieved from reporters table:", user.Email)
		return user, nil
	}

	logrus.Error("Error retrieving user:", err)
	return factory.User{}, err
}

func (p *Postgres) IsReporter(email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM reporters WHERE email = $1)`
	var exists bool
	err := p.dbConn.QueryRow(query, email).Scan(&exists)
	if err != nil {
		logrus.Error("Error checking if user is a reporter:", err)
		return false, err
	}
	logrus.Infof("Is reporter check for email %s: %t", email, exists)
	return exists, nil
}

func (p *Postgres) IncrementReporterReportCount(reporterID string) error {
	query := `UPDATE reporters SET total_reports = total_reports + 1 WHERE id = $1`
	_, err := p.dbConn.Exec(query, reporterID)
	if err != nil {
		logrus.Error("Error incrementing total_reports for reporter:", err)
		return err
	}
	logrus.Infof("Incremented report count for reporter ID: %s", reporterID)
	return nil
}
