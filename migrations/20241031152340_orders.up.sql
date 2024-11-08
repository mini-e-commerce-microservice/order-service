CREATE TABLE orders
(
    id                  bigserial primary key,
    user_id             bigint                   not null,
    status              varchar(100)             not null,
    total_amount        NUMERIC(10, 2)           not null,
    payment_status      varchar(100)             not null,
    payment_method_code varchar(100)             not null,
    tax                 NUMERIC(10, 2)           not null,
    shipping_cost       NUMERIC(10, 2)           not null,
    discount            NUMERIC(10, 2)           not null,
    order_date          timestamp with time zone not null,
    created_at          timestamp with time zone not null,
    updated_at          timestamp with time zone not null,
    deleted_at          timestamp with time zone
)