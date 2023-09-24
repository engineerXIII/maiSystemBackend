DROP TABLE IF EXISTS products CASCADE;

CREATE TABLE products
(
    product_id   UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    product_name VARCHAR(32)              NOT NULL CHECK ( product_name <> '' ),
    color        VARCHAR(32)              NOT NULL CHECK ( color <> '' ),
    description  VARCHAR(128)             NOT NULL CHECK ( description <> '' ),
    factory      VARCHAR(32)              NOT NULL CHECK ( factory <> '' ),
    cost         INTEGER                  NOT NULL CHECK ( cost > 0 ),
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE          DEFAULT CURRENT_TIMESTAMP
);