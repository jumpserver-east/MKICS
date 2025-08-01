package model

type MaintenanceRecord struct {
	ID                 int64  `json:"id"`
	ClientID           int    `json:"clientId"`
	MaintenanceTypes   string `json:"maintenanceTypes"`
	MaintenanceVersion string `json:"maintenanceVersion"`
	MaintenanceTitle   string `json:"maintenanceTitle"`
	MaintenanceContext string `json:"maintenanceContext"`
	MaintenanceTime    int64  `json:"maintenanceTime"`
	CreatedBy          string `json:"createdBy"`
	ModifiedAt         int64  `json:"modifiedAt"`
	ModifiedBy         string `json:"modifiedBy"`
	ModifiedByName     string `json:"modifiedByName"`
	CreatorName        string `json:"creatorName"`
	ClientName         string `json:"clientName"`
	CreateTime         int64  `json:"createTime"`
}
