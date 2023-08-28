-- +goose Up
-- +goose StatementBegin
CREATE TABLE container(
    name text PRIMARY KEY
);

CREATE TABLE container_mapping (
    long_name text PRIMARY KEY,
    container_name text,
    CONSTRAINT container_mapping_container_fk FOREIGN KEY (container_name) REFERENCES container(name) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    container_name text,
    log_text text,
    created_at INTEGER,
    CONSTRAINT log_container_fk FOREIGN KEY (container_name) REFERENCES container(name) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX idx_log_container ON log (container_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS log;
DROP TABLE IF EXISTS container_mapping;
DROP TABLE IF EXISTS container;
-- +goose StatementEnd
