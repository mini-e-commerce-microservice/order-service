CREATE TABLE product_items
(
    id                     bigserial primary key,
    user_id                BIGINT         NOT NULL,
    outlet_id              BIGINT         NOT NULL,
    variant_1              JSONB,
    variant_2              JSONB,
    sub_category_item_name VARCHAR(255),
    name                   VARCHAR(255)   NOT NULL,
    description            TEXT,
    price                  NUMERIC(15, 2) NOT NULL,
    stock                  INT            NOT NULL,
    sku                    VARCHAR(100),
    weight                 INT            NOT NULL,
    package_length         NUMERIC(10, 2),
    package_width          NUMERIC(10, 2),
    package_height         NUMERIC(10, 2),
    dimensional_weight     NUMERIC(10, 2),
    is_active              BOOLEAN                  DEFAULT TRUE,
    product_condition      VARCHAR(50),
    minimum_purchase       INT                      DEFAULT 1,
    size_guide_image       VARCHAR(255),
    trace_parent           varchar(255),
    created_at             TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at             TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at             TIMESTAMP WITH TIME ZONE
)