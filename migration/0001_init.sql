-- +goose Up
create extension if not exists "uuid-ossp";

create schema if not exists mart;

create table if not exists mart.usr(
    id uuid not null default uuid_generate_v4(),
    login varchar(64) not null,
    password varchar(128) not null,
    constraint usr_pkey primary key (id),
    constraint usr_login_uk unique (login)
);

create table if not exists mart.balance(
    id uuid not null default uuid_generate_v4(),
    cur numeric(10, 3) not null default 0,
    withdrawn numeric(10, 3) not null default 0,
    user_id uuid not null,
    constraint balance_pkey primary key (id),
    constraint cur check (cur >= 0),
    constraint balance_user_id_uk unique (user_id),
    constraint balance_user_id_fk foreign key(user_id) references mart.usr(id)
);

create table if not exists mart.ordr(
    id uuid not null default uuid_generate_v4(),
    num varchar(32) not null,
    status varchar(12) not null default 'NEW',
    accrual numeric(10, 3) default 0,
    user_id uuid not null,
    uploaded_at timestamp not null default CURRENT_TIMESTAMP,
    constraint ordr_pkey primary key (id),
    constraint ordr_num_uk unique (num)
);

create table if not exists mart.withdrawal(
    id uuid not null default uuid_generate_v4(),
    amount numeric(10, 3) not null,
    order_id uuid not null,
    processed_at timestamp not null default CURRENT_TIMESTAMP,
    constraint withdrawal_pkey primary key (id),
    constraint withdrawal_order_id_uk unique (order_id),
    constraint withdrawal_order_id_fk foreign key(order_id) references mart.ordr(id)
);
-- +goose Down