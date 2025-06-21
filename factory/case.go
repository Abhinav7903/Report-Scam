package factory

type Report struct {
	ID           string         `json:"id"`
	SubmittedBy  string         `json:"submitted_by"` // "user" or "reporter"
	SubmitterID  string         `json:"submitter_id"`
	ReportTypeID int            `json:"report_type_id"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	IsPublic     bool           `json:"is_public"`
	Status       string         `json:"status"`
	CreatedAt    string         `json:"created_at"`
	Metadata     []CaseMetadata `json:"metadata,omitempty"`
	Wallets      []Wallet       `json:"wallets,omitempty"`
	Domains      []Domain       `json:"domains,omitempty"`
}

type CaseMetadata struct {
	ID       string `json:"id"`
	ReportID string `json:"report_id"`
	Key      string `json:"key"`
	Value    string `json:"value"`
	IsPublic bool   `json:"is_public"`
}
type Wallet struct {
	ID       string `json:"id"`
	ReportID string `json:"report_id"`
	Address  string `json:"address"`
	Network  string `json:"network"`
}
type Domain struct { 
	ID         string `json:"id"`
	ReportID   string `json:"report_id"`
	DomainName string `json:"domain_name"`
}
type ReportType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
