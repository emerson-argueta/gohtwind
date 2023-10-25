package infra

type Pagination struct {
	TotalRecords int64 `json:"total_records"`
	TotalPages   int64 `json:"total_pages"`
	CurrentPage  int64 `json:"current_page"`
	PerPage      int64 `json:"per_page"`
}

func NewPagination(current_page int64, per_page int64, total_records int64) *Pagination {
	total_pages := total_records / per_page
	return &Pagination{
		TotalRecords: total_records,
		TotalPages:   total_pages,
		CurrentPage:  current_page,
		PerPage:      per_page,
	}
}
