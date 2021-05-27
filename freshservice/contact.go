package freshservice

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

// See here: https://api.freshservice.com/v2/#contact_attributes
type Contact struct {
	Active         bool                   `json:"active,omitempty"`
	Address        string                 `json:"address,omitempty"`
	CompanyId      int                    `json:"company_id,omitempty"`
	ViewAllTickets bool                   `json:"view_all_tickets,omitempty"`
	CustomFields   map[string]interface{} `json:"custom_fields,omitempty"`
	Deleted        bool                   `json:"deleted,omitempty"`
	Description    string                 `json:"description,omitempty"`
	Email          string                 `json:"email,omitempty"`
	Id             int                    `json:"id,omitempty"`
	JobTitle       string                 `json:"job_title,omitempty"`
	Language       string                 `json:"language,omitempty"`
	Mobile         string                 `json:"mobile,omitempty"`
	Name           string                 `json:"name,omitempty"`
	OtherEmails    []string               `json:"other_emails,omitempty"`
	Phone          string                 `json:"phone,omitempty"`
	Tags           []string               `json:"tags,omitempty"`
	Timezone       string                 `json:"time_zone,omitempty"`
	TwitterId      string                 `json:"twitter_id,omitempty"`
	ExternalId     string                 `json:"unique_external_id,omitempty"`
	OtherCompanies []string               `json:"other_companies,omitempty"`
	CreatedAt      time.Time              `json:"created_at,omitempty"`
	UpdatedAt      time.Time              `json:"updated_at,omitempty"`
}

const CONTACT_FILTER_LIMIT = 512

type PreparedContactQuery string
type ContactFilter struct {
	Field string
	Value interface{}
}

func formatValue(v interface{}) string {
	switch v.(type) {
	case int, int8, int16, int32, int64, float32, float64, bool:
		return fmt.Sprint(v)
	case string:
		return "'" + v.(string) + "'"
	default:
		return "'" + v.(string) + "'"
	}
}

// BuildQueryString returns a prepared query string that can be used to query
// Freshdesk with. The second parameter indicates if preparation was successful.
// If false, start again.
func BuildQueryString(query PreparedContactQuery, contact ContactFilter) (PreparedContactQuery, bool) {
	if len(query) == 0 {
		return PreparedContactQuery(fmt.Sprintf("\"%s: %s\"", contact.Field, formatValue(contact.Value))), true
	}

	prepared := fmt.Sprintf(" OR %s: %s", contact.Field, formatValue(contact.Value))
	if len(query)+len(prepared) <= CONTACT_FILTER_LIMIT {
		return PreparedContactQuery(string(query[:len(query)-1]) + prepared + "\""), true
	}

	return query, false
}

// FilterContacts filters contacts by some specified field (OR conditions only)
func (c *Client) FilterContacts(q PreparedContactQuery) ([]Contact, error) {
	var result struct {
		Total   int       `json:"total"`
		Results []Contact `json:"results"`
	}

	vars := url.Values{}
	vars.Set("query", string(q))
	uri := fmt.Sprintf("%s?%s", "/api/v2/search/contacts", vars.Encode())
	err := c.ReadObject(uri, &result)
	if err != nil {
		return nil, err
	}

	return result.Results, nil
}

// GetContacts returns a list of contacts.
// The API supports additional filters/adjustments that are not reflected here
// (to be implemented on an as-needed basis).
func (c *Client) GetContacts() ([]Contact, error) {
	var result []Contact

	err := c.ReadObject("/api/v2/contacts", &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) CreateContact(contact Contact) (*Contact, error) {
	if contact.Name == "" {
		return nil, errors.New("Name is a required field for new contacts")
	}
	if contact.Email == "" && contact.Phone == "" && contact.Mobile == "" && contact.TwitterId == "" && contact.ExternalId == "" {
		return nil, errors.New("At least one of: (email, phone, mobile, twitterid, externalid) must be provided")
	}

	var result *Contact
	err := c.WriteObject("/api/v2/contacts", "POST", contact, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateContact(contact Contact) (*Contact, error) {
	if contact.Id == 0 {
		return nil, errors.New("Need ID to update contact.")
	}

	var result *Contact
	uri := fmt.Sprintf("/api/v2/contacts/%d", contact.Id)
	err := c.WriteObject(uri, "PUT", contact, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
