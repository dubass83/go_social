CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level INT NOT NULL DEFAULT 0,
    description TEXT,
);

INSERT INTO roles (name, level, description) VALUES ('admin', 10, 'Administrator role');
INSERT INTO roles (name, level, description) VALUES ('moderator', 5, 'Moderator role');
INSERT INTO roles (name, level, description) VALUES ('user', 1, 'Default user role');
