package utils

import "strconv"

type SearchParams struct {
	Limit int
	Since string
	Decs bool
	Sort string
}

func (params *SearchParams) SearchParams () {
	params.Limit = -1;
	params.Since = "";
	params.Decs = false;
	params.Sort = "";
}

func (params *SearchParams) CreateParams (limit string, since string, desc string, sort string) {
	params.SearchParams()
	if limit != "" {
		params.Limit, _ = strconv.Atoi(limit)
	}
	params.Since = since
	if desc == "true" {
		params.Decs = true
	}
}

