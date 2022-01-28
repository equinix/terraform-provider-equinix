package api

//Pagination describes pagination controls for Network Edge device type query
type Pagination struct {
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
	Total  int `json:"total,omitempty"`
}
