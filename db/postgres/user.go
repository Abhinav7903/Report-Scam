package postgres

import (
	"abuse/factory"
	"fmt"

	"github.com/sirupsen/logrus"
)

func (p *Postgres) CreateUser(data factory.User) (string, error) {
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

func (p *Postgres) GetUser(email string) (factory.User, error) {
	query := `SELECT id, email, name, contact_info, affiliation, total_reports, user_type, created_at
			  FROM users WHERE email = $1
			  UNION ALL
			  SELECT id, email, name, contact_info, affiliation, total_reports, user_type, created_at
			  FROM reporters WHERE email = $1`
	row := p.dbConn.QueryRow(query, email)

	var user factory.User
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.ContactInfo,
		&user.Affiliation, &user.TotalReports, &user.UserType, &user.CreatedAt)
	if err != nil {
		logrus.Error("Error retrieving user:", err)
		return factory.User{}, err
	}

	logrus.Info("User retrieved successfully:", user.Email)
	return user, nil
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
