package log

type SearchRequest struct {
	Contains string `json:"contains"`
}

type GetLogsResponse struct {
	Records     []string `json:"records"`
	FirstCursor int64    `json:"firstCursor"`
	LastCursor  int64    `json:"lastCursor"`
}
