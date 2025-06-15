-- name: CreateUser :exec
INSERT INTO users (zitadel_id,
                   first_name,
                   last_name,
                   username,
                   email)
VALUES ($1, $2, $3, $4, $5);
