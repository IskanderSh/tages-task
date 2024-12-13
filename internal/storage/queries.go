package storage

const (
	queryGetByID = `
SELECT id, file_name, path, created_at, updated_at
FROM file_meta
WHERE file_name LIKE CONCAT(SPLIT_PART($1, '.', 1), '%.', SPLIT_PART($1, '.', 2))
ORDER BY created_at DESC
LIMIT 1;`

	queryCreate = `
INSERT INTO file_meta (file_name, path, created_at, updated_at)
VALUES ($1, $2, $3, $4);
`

	queryFetchAll = `
SELECT id, file_name, path, created_at, updated_at
FROM file_meta
ORDER BY created_at DESC;
`
)
