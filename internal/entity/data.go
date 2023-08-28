package entity

import "fmt"

type Direction string

const (
	DirFuture = "future"
	DirPast   = "past"
)

type GetLogsRequest struct {
	ShortName   string
	Direction   Direction `form:"direction"`
	ChunkNumber int       `form:"chunk_number"`
	Position    int       `form:"position"`
	Limit       int       `form:"limit"`
}

func (r *GetLogsRequest) IsValid() error {
	if r.Limit == 0 {
		return fmt.Errorf("limit should be positive")
	}
	return nil
}

type ContainerInfo struct {
	ShortName string `json:"shortName"`
	IsAlive   bool   `json:"isAlive"`
}
