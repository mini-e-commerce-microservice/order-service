CREATE TABLE outbox_events
(
    id            uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    aggregatetype varchar(255) not null,
    aggregateid   varchar(255) not null,
    type          varchar(255) not null,
    payload       jsonb        not null,
    trace_parent  varchar(255)
)