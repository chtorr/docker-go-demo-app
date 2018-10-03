CREATE TABLE demo (
    id BIGSERIAL NOT NULL,
    name TEXT,

    CONSTRAINT demo_pk1 PRIMARY KEY(id)
);

INSERT INTO demo (name) VALUES ('test name');
