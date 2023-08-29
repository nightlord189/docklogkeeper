package log

import (
	"github.com/nightlord189/docklogkeeper/internal/entity"
)

func (a *Adapter) GetLogs(req entity.GetLogsRequest) (GetLogsResponse, error) {
	logs, err := a.Repo.GetLogs(req.ShortName, req.Direction == entity.DirFuture, req.Cursor, req.Limit)
	if err != nil {
		return GetLogsResponse{}, err
	}

	if len(logs) == 0 {
		return GetLogsResponse{Records: []string{}}, nil
	}

	lines := make([]string, len(logs))
	for i := range logs {
		lines[i] = logs[i].LogText
	}

	return GetLogsResponse{
		Records:     lines,
		FirstCursor: logs[len(logs)-1].ID, //earlier
		LastCursor:  logs[0].ID,           //later
	}, nil
}
