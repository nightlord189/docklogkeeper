package log

type SearchRequest struct {
	Contains string `json:"contains"`
}

type GetLinesRequest struct {
	ShortName   string `json:"shortName"`
	ChunkNumber int    `json:"chunkNumber"`
	Offset      int    `json:"offset"`
	Limit       int    `json:"limit"`
}

type GetLinesResponse struct {
	Records     []string `json:"records"`
	ChunkNumber int      `json:"chunkNumber"`
	Offset      int      `json:"offset"`
}
