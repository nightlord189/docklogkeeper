package trigger

import (
	"net/http"
	"sync"
	"time"

	"github.com/nightlord189/docklogkeeper/internal/entity"
	"github.com/nightlord189/docklogkeeper/internal/repo"
)

const bufferSize = 1000

type Adapter struct {
	LogsChan      chan entity.LogDataDB
	Repo          *repo.Repo
	triggersCache *TriggersCache
	regexpCache   *sync.Map // regexp string -> regexp.Regexp
	httpClient    *http.Client
}

func New(repoInst *repo.Repo) *Adapter {
	return &Adapter{
		Repo:          repoInst,
		LogsChan:      make(chan entity.LogDataDB, bufferSize),
		triggersCache: NewTriggersCache(),
		regexpCache:   &sync.Map{},
		httpClient: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       1 * time.Minute,
		},
	}
}
