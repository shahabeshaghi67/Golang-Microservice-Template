CREATE TABLE users
(
    id           UUID,
    name         VARCHAR(100) NOT NULL,
    created_at   TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP    NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);
