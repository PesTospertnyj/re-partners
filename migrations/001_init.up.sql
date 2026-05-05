CREATE TABLE IF NOT EXISTS pack_sizes
(
    id   BIGSERIAL PRIMARY KEY,
    size INTEGER CHECK (size > 0) UNIQUE
);

INSERT INTO pack_sizes (size)
VALUES (250),
       (500),
       (1000),
       (2000),
       (5000)
ON CONFLICT (size) DO NOTHING;
