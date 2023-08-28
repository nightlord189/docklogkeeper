package handler

import (
	"fmt"
	"github.com/nightlord189/docklogkeeper/internal/entity"
)

type SearchLogsRequest struct {
	Contains string `form:"contains"`
}

type SearchLogsResponse struct {
	Records []string `json:"records"`
}

type GetContainersResponse struct {
	Containers []entity.ContainerInfo `json:"containers"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GenericResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func GenericError(message string) GenericResponse {
	return GenericResponse{
		Message: message,
	}
}

func GenericErrorf(message string, args ...any) GenericResponse {
	return GenericResponse{
		Message: fmt.Sprintf(message, args...),
	}
}
