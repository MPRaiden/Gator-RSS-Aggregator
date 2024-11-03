-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
    )
RETURNING *;

-- name: ListFeeds :many
select feeds.name, feeds.url, users.name
from feeds
join users on feeds.user_id = users.id
order by feeds.name
;

-- name: GetFeedByURL :one
select *
from feeds
where url = $1
;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

