CREATE TABLE IF NOT EXISTS stacks
(
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,

    ulid       VARCHAR(26)  NOT NULL UNIQUE,
    name        VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL,
    account_id  BIGINT       NOT NULL,
    CONSTRAINT fk_stacks_account_id FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE ON UPDATE CASCADE
);
