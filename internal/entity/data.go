package entity

import "fmt"

type Direction string

const (
	DirFuture = "future"
	DirPast   = "past"
)

type GetLogsRequest struct {
	ShortName string
	Direction Direction `form:"direction"`
	Cursor    int64     `form:"cursor"`
	Limit     int       `form:"limit"`
}

func (r *GetLogsRequest) IsValid() error {
	if r.Limit == 0 {
		return fmt.Errorf("limit should be positive")
	}
	if r.Direction != DirFuture && r.Direction != DirPast {
		return fmt.Errorf("invalid direction")
	}
	return nil
}

type ContainerInfo struct {
	ShortName string `json:"shortName"`
	IsAlive   bool   `json:"isAlive"`
}
