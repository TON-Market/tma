package db

var q = `
create table if not exists assets
(
    user_raw_address  varchar(255) not null,
    event_id          uuid         not null,
    collateral_staked bigint       not null,
    token             varchar(50)  not null,
    size              bigint       not null,
    primary key (user_raw_address, event_id, token)
);

create table if not exists users
(
    raw_addr varchar(255) not null
        primary key
);

create table if not exists deals
(
    id            uuid              not null
        primary key,
    event_id      uuid              not null,
    token         varchar(10)       not null,
    collateral    bigint            not null,
    size          bigint            not null,
    user_raw_addr varchar(255)
        constraint fk_user_raw_addr
            references users
            on delete cascade,
    deal_status   integer default 0 not null,
	attempts integer default 0
);

create table if not exists user_deals
(
    user_raw_addr varchar(255)
        references users
            on delete cascade,
    deal_id       uuid
        references deals


            on delete cascade
);
`
