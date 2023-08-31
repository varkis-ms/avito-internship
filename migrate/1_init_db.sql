CREATE TABLE IF NOT EXISTS Segments
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR   NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    deleted_at    timestamptz          DEFAULT NULL,
    UNIQUE (name)
);


CREATE TABLE IF NOT EXISTS Users
(
    id INTEGER PRIMARY KEY
);


CREATE TABLE IF NOT EXISTS Users_segment
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    INTEGER   NOT NULL REFERENCES Users (id),
    segment_id INTEGER   not null REFERENCES Segments (id),
    added_at timestamptz NOT NULL DEFAULT now(),
    left_at    timestamptz          DEFAULT NULL
);

CREATE INDEX ON Users_segment (user_id);
