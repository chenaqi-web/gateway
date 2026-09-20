package vaildate

func Page(page int) int {
	if page < 1 {
		page = 1
	}
	return page
}

func PageSize(pageSize int) int {
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}
	return pageSize
}
