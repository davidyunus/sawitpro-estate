CREATE TABLE estate (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(36) NOT NULL,
    length INT NOT NULL,
    width INT NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT NOW()
)

CREATE TABLE palmTreeLocation (
    id SERIAL PRIMARY KEY,
    estateId BIGINT NOT NULL,
    uuid VARCHAR(36) NOT NULL,
    x INT NOT NULL,
    y INT NOT NULL,
    height INT NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT NOW()
)