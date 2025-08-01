package model

type Subscription struct {
	ID                   int64          `json:"id"`
	ClientID             int            `json:"clientId"`
	ProductServiceID     int            `json:"productServiceId"`
	Client               Client         `json:"client"`
	ProductService       ProductService `json:"productService"`
	SubscriptionType     string         `json:"subscriptionType"`
	Expired              bool           `json:"expired"`
	SalesUserID          string         `json:"salesUserId"`
	SalesTeam            string         `json:"salesTeam"`
	SalesUser            SalesUser      `json:"salesUser"`
	Amount               int            `json:"amount"`
	StartDate            int64          `json:"startDate"`
	EndDate              int64          `json:"endDate"`
	SupportEndDate       int64          `json:"supportEndDate"`
	SupportExpired       bool           `json:"supportExpired"`
	ServiceType          string         `json:"serviceType"`
	Description          string         `json:"description"`
	Region               string         `json:"region"`
	AcceptanceReport     bool           `json:"acceptanceReport"`
	Internal             bool           `json:"internal"`
	MarketingGrade       string         `json:"marketingGrade"`
	EnterprisePreOrder   bool           `json:"enterprisePreOrder"`
	SerialNo             string         `json:"serialNo"`
	SubscriptionTypeName string         `json:"subscriptionTypeName"`
	ServiceTypeName      string         `json:"serviceTypeName"`
}

type ProductService struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Amount          int    `json:"amount"`
	ProductID       int    `json:"productId"`
	UpdateTimestamp int64  `json:"updateTimestamp"`
	Deadline        int64  `json:"deadline"`
	AmountUnit      string `json:"amountUnit"`
	MaxPLU          int    `json:"maxPlu"`
	CreateTimestamp int64  `json:"createTimestamp"`
	Enabled         bool   `json:"enabled"`
}
