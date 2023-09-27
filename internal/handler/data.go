package handler

import (
	"fmt"
	"github.com/nightlord189/docklogkeeper/internal/entity"
)

type TemplateData struct {
	Analytics bool
}

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

type CreateTriggerRequest struct {
	Name           string `json:"name"`
	ContainerName  string `json:"containerName"`
	Contains       string `json:"contains"`
	NotContains    string `json:"notContains"`
	Regexp         string `json:"regexp"`
	Method         string `json:"method"`
	WebhookURL     string `json:"webhookURL"`
	WebhookHeaders string `json:"webhookHeaders"`
	WebhookBody    string `json:"webhookBody"`
}

func (r *CreateTriggerRequest) ToDB() entity.TriggerDB {
	return entity.TriggerDB{
		ID:             0,
		Name:           r.Name,
		ContainerName:  r.ContainerName,
		Contains:       r.Contains,
		NotContains:    r.NotContains,
		Regexp:         r.Regexp,
		Method:         r.Method,
		WebhookURL:     r.WebhookURL,
		WebhookHeaders: r.WebhookHeaders,
		WebhookBody:    r.WebhookBody,
	}
}

type GetTriggersResponse struct {
	Records []entity.TriggerDB `json:"records"`
}
