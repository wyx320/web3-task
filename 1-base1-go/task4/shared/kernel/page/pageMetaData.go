package page

import "math"

type PageMetaData struct {
	PageIndex  int `json:"page_index"`
	PageSize   int `json:"page_size"`
	TotalPage  int `json:"total"`
	TotalCount int `json:"total_count"`

	HasPrevious bool `json:"has_previous"`
	HasNext     bool `json:"has_next"`
}

func NewPageMetaData(totalCount int, pageObject PageObject) *PageMetaData {
	total := int(totalCount)
	pageIndex := pageObject.PageIndex
	pageSize := pageObject.PageSize

	return &PageMetaData{
		PageIndex:   pageIndex,
		PageSize:    pageSize,
		TotalPage:   int(math.Ceil(float64(totalCount) / float64(pageSize))),
		TotalCount:  totalCount,
		HasPrevious: pageIndex > 1,
		HasNext:     pageIndex < total,
	}
}
