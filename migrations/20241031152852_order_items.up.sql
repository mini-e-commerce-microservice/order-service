CREATE TABLE order_items
(
    id                 bigserial primary key,
    order_id           bigint references orders (id) ON DELETE cascade,
    product_item_id    bigint references product_items (id) ON DELETE cascade,
    name               VARCHAR(255)   NOT NULL,
    description        TEXT,
    weight             INT            NOT NULL,
    package_length     NUMERIC(10, 2),
    package_width      NUMERIC(10, 2),
    package_height     NUMERIC(10, 2),
    dimensional_weight NUMERIC(10, 2),
    qty                int            not null,
    unit_price         NUMERIC(10, 2) not null,
    total_price        NUMERIC(10, 2) not null,
    discount           NUMERIC(10, 2) not null
)