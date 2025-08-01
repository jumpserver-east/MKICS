package model

type PaginationResponse[T any] struct {
	Total int `json:"marker"`
	Data  T   `json:"data"`
}
