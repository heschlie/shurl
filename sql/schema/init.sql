CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  username text NOT NULL
);

CREATE TABLE shurls (
  id BIGSERIAL PRIMARY KEY,
  hash text NOT NULL,
  url text NOT NULL,
  owner BIGINT REFERENCES users(id),
  hits BIGINT NOT NULL DEFAULT 0,
  expire BIGINT NOT NULL DEFAULT 90
);
