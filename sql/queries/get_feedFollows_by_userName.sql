-- name: GetFeedFollowsForUser :many
WITH feed_follow AS (
    SELECT * FROM feed_follows
)
SELECT
    feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM feed_follow
INNER JOIN feeds ON feed_follow.feed_id = feeds.id
INNER JOIN users ON feed_follow.user_id = users.id AND users.name = $1;
