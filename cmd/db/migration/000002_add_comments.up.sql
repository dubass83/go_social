CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGSERIAL NOT NULL REFERENCES users(id),
    post_id BIGSERIAL NOT NULL REFERENCES posts(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW()
);
