package factory

type User struct {
	ID           string `json:"id"` // UUID
	Email        string `json:"email"`
	Name         string `json:"name"`
	ContactInfo  string `json:"contact_info,omitempty"` // for users
	Affiliation  string `json:"affiliation,omitempty"`  // for reporters
	TotalReports int    `json:"total_reports,omitempty"`
	UserType     string `json:"user_type"` // "user" or "reporter"
	CreatedAt    string `json:"created_at"`
}
