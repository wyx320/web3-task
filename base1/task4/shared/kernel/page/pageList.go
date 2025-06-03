package page

type PageList[T any] struct {
	Items    []T           `json:"items"`
	MetaData *PageMetaData `json:"meta_data"`
}

func NewPageList[T any](items []T, totalCount int, pageObject PageObject) *PageList[T] {
	pageMetaData := NewPageMetaData(totalCount, pageObject)
	return &PageList[T]{
		Items:    items,
		MetaData: pageMetaData,
	}
}
