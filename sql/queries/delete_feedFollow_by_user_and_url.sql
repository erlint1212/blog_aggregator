-- name: DeleteFeedFollowByUrlAndName :exec
DELETE FROM feed_follows
WHERE feed_id IN (SELECT f.id FROM feeds f WHERE f.url = $1)
AND user_id IN (SELECT u.id FROM users u WHERE u.name = $2);
