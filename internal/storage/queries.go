package storage

const (
	queryGetByID = `
SELECT id, file_name, path, created_at, updated_at
FROM file_meta
WHERE file_name LIKE CONCAT($1::text, '%')
ORDER BY created_at DESC
LIMIT 1;`

	queryCreate = `
INSERT INTO file_meta (file_name, path, created_at, updated_at)
VALUES ($1, $2, $3, $4);
`
)
