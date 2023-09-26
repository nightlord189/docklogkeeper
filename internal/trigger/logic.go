package trigger

import (
	"bytes"
	"context"
	"github.com/nightlord189/docklogkeeper/internal/entity"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

func (a *Adapter) ReloadCache(ctx context.Context) {
	triggers, err := a.Repo.GetAllTriggers()
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("ReloadCache: get triggers error")
		return
	}

	a.triggersCache.Clear()
	for _, trig := range triggers {
		var currentArray []entity.TriggerDB
		gotValue, ok := a.triggersCache.Load(trig.ContainerName)
		if !ok {
			currentArray = make([]entity.TriggerDB, 0, 3)
		} else {
			currentArray = gotValue
		}
		currentArray = append(currentArray, trig)
		a.triggersCache.Store(trig.ContainerName, currentArray)
	}

	clearSyncMap(a.regexpCache)
	for _, trig := range triggers {
		if trig.Regexp == "" {
			continue
		}
		_, ok := a.regexpCache.Load(trig.Regexp)
		if !ok {
			newRegexp, err := regexp.Compile(trig.Regexp)
			if err != nil {
				log.Ctx(ctx).Err(err).Str("regexp", trig.Regexp).Msg("ReloadCache: compile regexp error")
				continue
			}
			a.regexpCache.Store(trig.Regexp, newRegexp)
		}
	}

	log.Ctx(ctx).Debug().Msg("triggers cache reloaded")
}

const workersCount = 3

func (a *Adapter) Run(ctx context.Context) {
	a.ReloadCache(ctx)

	wg := &sync.WaitGroup{}
	wg.Add(workersCount)
	for i := 0; i < workersCount; i++ {
		go a.readInput(ctx, wg)
	}
	wg.Wait()
}

func (a *Adapter) readInput(ctx context.Context, wg *sync.WaitGroup) {
	for logEntry := range a.LogsChan {
		gotTriggers := a.triggersCache.LoadWithAll(logEntry.ContainerName)
		if len(gotTriggers) == 0 {
			continue
		}
		a.processTriggers(ctx, &logEntry, gotTriggers)
	}
	wg.Done()
}

func (a *Adapter) processTriggers(ctx context.Context, logEntry *entity.LogDataDB, triggers []entity.TriggerDB) {
	for _, trig := range triggers {
		if trig.Match(logEntry.LogText, a.getRegexpFromCache(trig.Regexp)) {
			log.Ctx(ctx).Info().Msgf("trigger [%d %s] matched with log %s", trig.ID, trig.Name, logEntry.LogText)
			go a.sendWebhook(ctx, logEntry, &trig)
		}
	}
}

func (a *Adapter) sendWebhook(ctx context.Context, logEntry *entity.LogDataDB, trigger *entity.TriggerDB) {
	rawURL := injectVariables(trigger.WebhookURL, logEntry)
	rawBody := injectVariables(trigger.WebhookBody, logEntry)
	rawHeaders := trigger.WebhookHeaders
	if rawHeaders == "" {
		rawHeaders = entity.DefaultHeaders
	}
	rawHeaders = injectVariables(trigger.WebhookHeaders, logEntry)

	var body io.Reader
	if rawBody != "" {
		body = bytes.NewBuffer([]byte(rawBody))
	}

	httpReq, err := http.NewRequestWithContext(ctx, rawURL, rawURL, body)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("sendWebhook: create http request error")
		return
	}

	for header, value := range convertHeaders(rawHeaders) {
		httpReq.Header.Set(header, value)
	}

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("sendWebhook: http error")
		return
	}

	defer resp.Body.Close()

	readResp, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("sendWebhook: read response error")
	}

	log.Ctx(ctx).Debug().Msgf("webhook sent: [url %s], response: %s", rawURL, string(readResp))
}

func convertHeaders(raw string) map[string]string {
	result := make(map[string]string)
	splitted := strings.Split(raw, ";")
	for _, header := range splitted {
		splittedHeader := strings.Split(header, ":")
		if len(splittedHeader) != 2 {
			continue
		}
		result[splittedHeader[0]] = splittedHeader[1]
	}
	return result
}

func injectVariables(initialStr string, logEntry *entity.LogDataDB) string {
	if initialStr == "" {
		return initialStr
	}
	replacer := strings.NewReplacer(
		"$dlk_container_full_name", logEntry.ContainerFullName,
		"$dlk_container_name", logEntry.ContainerName,
		"$dlk_log", logEntry.LogText,
		"$dlk_timestamp", time.Unix(logEntry.CreatedAt, 0).Format(time.RFC3339),
	)
	return replacer.Replace(initialStr)
}
