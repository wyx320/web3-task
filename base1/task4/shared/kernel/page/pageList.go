package page

type PageList[T any] struct {
	Items    []T           `json:"items"`
	MetaData *PageMetaData `json:"meta_data"`
}

func NewPageList[T any](items []T, pageObject PageObject) *PageList[T] {
	pageMetaData := NewPageMetaData(len(items), pageObject)
	return &PageList[T]{
		Items:    items,
		MetaData: pageMetaData,
	}
}
