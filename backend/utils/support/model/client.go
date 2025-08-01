package model

type Client struct {
	ID                int64         `json:"id"`
	Name              string        `json:"name"`
	AbbreviatedName   string        `json:"abbreviatedName"`
	Address           string        `json:"address"`
	Description       string        `json:"description"`
	CreatedTimestamp  int64         `json:"createdTimestamp"`
	SalesUsers        []SalesUser   `json:"salesUsers"`
	SalesUserIds      []string      `json:"salesUserIds"`
	CustomerChurnType int           `json:"customerChurnType"`
	ClientParent      int           `json:"clientParent"`
	ClientType        string        `json:"clientType"`
	SupportUserIds    []string      `json:"supportUserIds"`
	SupportUsers      []SupportUser `json:"supportUsers"`
	NameValid         bool          `json:"nameValid"`
	ClientCode        string        `json:"clientCode"`
}

type SupportUser struct {
	UserId                string   `json:"userId"`
	Username              string   `json:"username"`
	CreateTime            int64    `json:"createTime"`
	Email                 string   `json:"email"`
	Phone                 string   `json:"phone"`
	Name                  string   `json:"name"`
	GroupId               int      `json:"groupId"`
	Active                bool     `json:"active"`
	ProductIds            []int    `json:"productIds"`
	RegionId              int      `json:"regionId"`
	Region                Region   `json:"region"`
	Group                 Group    `json:"group"`
	GroupType             string   `json:"groupType"`
	RegionType            string   `json:"regionType"`
	SupportedProductTypes []string `json:"supportedProductTypes"`
}

type Region struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	RegionType string `json:"regionType"`
}

type Group struct {
	EnName                   string      `json:"enName"`
	ID                       int         `json:"id"`
	Name                     string      `json:"name"`
	KcGroupId                interface{} `json:"kcGroupId"`
	EnableRegion             bool        `json:"enableRegion"`
	EnableProduct            bool        `json:"enableProduct"`
	LeastManagementPrivilege string      `json:"leastManagementPrivilege"`
}
