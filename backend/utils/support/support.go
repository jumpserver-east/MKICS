package support

import (
	"EvoBot/backend/utils/support/httplib"
	"EvoBot/backend/utils/support/model"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func NewSupportClient(baseUrl string, username string, password string) (*Client, error) {
	client, err := httplib.NewClient(baseUrl, 30*time.Second)
	if err != nil {
		return nil, err
	}
	client.SetBasicAuth(username, password)
	return &Client{
		BaseURL: baseUrl,
		client:  client,
	}, nil
}

type Client struct {
	BaseURL string
	client  *httplib.Client
}

type BearerTokenAuth struct {
	Token string
}

func (auth *BearerTokenAuth) Sign(r *http.Request) error {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", auth.Token))
	return nil
}

func (c *Client) GetClients(marker int, max int) (ret model.PaginationResponse[[]model.Client], err error) {
	_, err = c.client.Get(SupportClientsURL, &ret, map[string]string{
		"region": "eastern",
		"simple": "false",
		"marker": strconv.Itoa(marker),
		"max":    strconv.Itoa(max),
	})
	return
}

func (c *Client) GetMaintenances(marker int, max int) (ret model.PaginationResponse[[]model.Maintenance], err error) {
	_, err = c.client.Get(SupportMaintenancesURL, &ret, map[string]string{
		"product": "JumpServer",
		"region":  "eastern",
		"marker":  strconv.Itoa(marker),
		"max":     strconv.Itoa(max),
	})
	return
}

func (c *Client) GetMaintenanceRecords(marker int, max int) (ret model.PaginationResponse[[]model.MaintenanceRecord], err error) {
	_, err = c.client.Get(SupportMaintenanceRecordsURL, &ret, map[string]string{
		"region": "eastern",
		"simple": "false",
		"marker": strconv.Itoa(marker),
		"max":    strconv.Itoa(max),
	})
	return
}

func (c *Client) GetSubscriptionsWithQuickSearch(marker int, max int, quickSearch string) (ret model.PaginationResponse[[]model.Subscription], err error) {
	_, err = c.client.Get(SupportSubscriptionsURL, &ret, map[string]string{
		"product":          "JumpServer",
		"region":           "eastern",
		"subscriptionType": "enterprise",
		"simple":           "false",
		"quickSearch":      quickSearch,
		"marker":           strconv.Itoa(marker),
		"max":              strconv.Itoa(max),
	})
	return
}
