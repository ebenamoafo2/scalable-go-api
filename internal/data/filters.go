package data

import "github.com/ebenamoafo2/scalable-go-api/internal/validator"

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	//Check that the page and the page size are valid
	v.Check(f.Page > 0, "page", "must be greater than  0")
	v.Check(f.Page <= 10_000_000, "page_size", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than 0")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

	v.Check(validator.PermittedValue(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}
