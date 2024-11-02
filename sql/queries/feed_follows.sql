-- name: CreateFeedFollow :many
with
    inserted_feed_follow as (
        insert into feed_follows(id, created_at, updated_at, user_id, feed_id)
        values ($1, $2, $3, $4, $5)
        returning *
    )
select inserted_feed_follow.*, feeds.name as feed_name, users.name as user_name
from inserted_feed_follow
inner join feeds on feeds.id = inserted_feed_follow.feed_id
inner join users on users.id = inserted_feed_follow.user_id
;

-- name: GetFeedFollowsForUser :many
select feed_follows.*, feeds.name as feed_name, users.name as user_name
from feed_follows
inner join feeds on feed_follows.feed_id = feeds.id
inner join users on feed_follows.user_id = users.id
where feed_follows.user_id = $1
;

-- name: DeleteFeedFollow :exec
delete from feed_follows
where
    feed_follows.user_id = $1
    and feed_follows.feed_id = (select id from feeds where feeds.url = $2)
;

