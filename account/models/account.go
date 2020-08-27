package models

// Account - represents an account object
type Account struct {
	Attributes     Attributes `json:"attributes"`
	CreatedOn      string     `json:"created_on,omitempty"`
	ID             string     `json:"id,omitempty"`
	ModifiedOn     string     `json:"modified_on,omitempty"`
	OrganisationID string     `json:"organisation_id,omitempty"`
	Type           string     `json:"type,omitempty"`
	Version        int32      `json:"version,omitempty"`
}
