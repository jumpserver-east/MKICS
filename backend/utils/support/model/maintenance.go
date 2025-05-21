package model

type Maintenance struct {
	ID               int64   `json:"id"`
	SubscriptionID   int     `json:"subscriptionId"`
	Status           string  `json:"status"`
	Content          Content `json:"content"`
	DeploymentTime   int64   `json:"deploymentTime"`
	DeploymentMethod string  `json:"deploymentMethod"`
	Template         string  `json:"template"`
	CreatedBy        string  `json:"createdBy"`
	CreatorName      string  `json:"creatorName"`
	ModifiedAt       int64   `json:"modifiedAt"`
	ModifiedBy       string  `json:"modifiedBy"`
	ModifiedByName   string  `json:"modifiedByName"`
	CreateTime       int64   `json:"createTime"`
}

type Content struct {
	Elements []Element `json:"elements"`
}

type Element struct {
	Title      string            `json:"title"`
	ContentMap map[string]string `json:"contentMap"`
}
