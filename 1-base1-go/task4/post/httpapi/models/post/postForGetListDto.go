package models

type PostForGetListDto struct {
	Title  string `json:"title"`
	UserId uint64 `json:"user_id"`

	PageSize  int `json:"page_size"`
	PageIndex int `json:"page_index"`
	// total     int64 `json:"total"`
}
