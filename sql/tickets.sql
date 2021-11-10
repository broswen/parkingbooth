CREATE TABLE tickets (
    id text,
    location_id text,
    start_epoch bigint,
    stop_epoch bigint,
    duration_seconds bigint,
    payment_id text,
    UNIQUE (id, location_id)
);