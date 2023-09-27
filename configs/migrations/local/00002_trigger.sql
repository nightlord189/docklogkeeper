-- +goose Up
-- +goose StatementBegin

-- internal variables:
-- $dlk_container_full_name
-- $dlk_container_name
-- $dlk_log
-- $dlk_timestamp

CREATE TABLE trigger (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    trigger_name text,
    container_name text,
    contains text NULL,
    not_contains text NULL,
    regexp text NULL,
    method text, --webhook
    webhook_url text NULL, -- formatted url for webhook
    webhook_headers text NULL, -- formatted list of headers for webhook delimited by ; example: Content-Type:application/json;Accept:*/*;Authorization:token
    webhook_body text NULL -- formatted json for webhook
);

CREATE UNIQUE INDEX idx_trigger_name ON trigger (trigger_name);

ALTER TABLE log ADD container_full_name text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE log DROP container_full_name;

DROP INDEX IF EXISTS idx_trigger_name;
DROP TABLE IF EXISTS trigger;
-- +goose StatementEnd
