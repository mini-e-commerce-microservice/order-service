CREATE TABLE saga_states
(
    id      bigserial PRIMARY KEY,
    payload jsonb        not null,
    status  varchar(100) not null,
    step    jsonb        not null,
    type    varchar(100) not null,
    version varchar(10)  not null
)