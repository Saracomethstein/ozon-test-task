CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    parent_id BIGINT REFERENCES comments(id) ON DELETE CASCADE,
    author TEXT NOT NULL,
    body TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT to_char(now(), 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"')
);

CREATE INDEX IF NOT EXISTS idx_comments_post_created_id ON comments (post_id, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_comments_parent_created_id ON comments (parent_id, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);