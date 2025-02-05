-- name: GetPostsForUser :many
SELECT p.* 
FROM posts p
INNER JOIN feeds f ON p.feed_id = f.id
INNER JOIN users u ON u.id = f.user_id AND u.name = $1
ORDER BY p.published_at ASC
LIMIT $2;
