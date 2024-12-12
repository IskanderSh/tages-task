CREATE TABLE IF NOT EXISTS file_meta (
    id serial,
    file_name text unique,
    path text,
    created_at timestamp,
    updated_at timestamp
);