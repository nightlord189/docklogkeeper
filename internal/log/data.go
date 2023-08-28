package log

type SearchRequest struct {
	Contains string `json:"contains"`
}

type GetLogsResponse struct {
	Records     []string `json:"records"`
	ChunkNumber int      `json:"chunkNumber"`
	Offset      int      `json:"offset"`
}
