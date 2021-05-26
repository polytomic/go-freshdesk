package freshservice

import (
	"fmt"
	"net/url"
	"time"
)

//Ticket see here: https://api.freshservice.com/v2/#contact_attributes
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

// BuildQueryString takes an array of FreshdeskContactQuery into a maximally-
// sized query string, returning that and also a slice of the remaining objects.
func BuildQueryString(query []ContactFilter) (PreparedContactQuery, []ContactFilter) {
	if len(query) < 1 {
		return "", query
	}

	qs := fmt.Sprintf("%s: %s", query[0].Field, formatValue(query[0].Value))
	numTaken := 1

	for _, q := range query {
		prepared := fmt.Sprintf(" OR %s: %s", q.Field, formatValue(q.Value))

		if len(qs)+len(prepared) > CONTACT_FILTER_LIMIT-2 {
			break
		}

		qs += prepared
		numTaken += 1
	}

	return PreparedContactQuery(fmt.Sprintf("\"%s\"", qs)), query[numTaken:]
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
