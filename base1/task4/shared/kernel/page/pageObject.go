package page

type PageObject struct {
	PageSize  int `json:"page_size"`
	PageIndex int `json:"page_index"`
}

const maxPageSize = 100

func NewPageObject(pageSize, pageIndex int) *PageObject {
	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize <= 0 || pageSize > maxPageSize {
		pageSize = 50 // 默认分页大小
	}

	return &PageObject{
		PageSize:  pageSize,
		PageIndex: pageIndex,
	}
}
