package model

type SalesUser struct {
	UserId      string `json:"userId"`
	Username    string `json:"username"`
	CreateTime  int64  `json:"createTime"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Name        string `json:"name"`
	GroupId     int    `json:"groupId"`
	Active      bool   `json:"active"`
	ProductIds  []int  `json:"productIds"`
	RegionId    int    `json:"regionId"`
	SalesTeamId int    `json:"salesTeamId"`
	Region      Region `json:"region"`
	Group       Group  `json:"group"`
	GroupType   string `json:"groupType"`
	RegionType  string `json:"regionType"`
}

type SalesTeam struct {
	Label           string `json:"label"`
	UpdateTimestamp int64  `json:"updateTimestamp"`
	CreateTimestamp int64  `json:"createTimestamp"`
	LeaderId        string `json:"leaderId"`
	Name            string `json:"name"`
	Parent          int    `json:"parent"`
	ID              int    `json:"id"`
}
