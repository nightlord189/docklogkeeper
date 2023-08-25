package log

type SearchRequest struct {
	Contains string `json:"contains"`
}

type GetLinesRequest struct {
	ShortName  string `json:"shortName"`
	FileNumber int    `json:"fileNumber"`
	Offset     int    `json:"offset"`
	Limit      int    `json:"limit"`
}

type GetLinesResponse struct {
	Lines      []string `json:"lines"`
	FileNumber int      `json:"fileNumber"`
	Offset     int      `json:"offset"`
}
