-- +goose Up
create extension if not exists "uuid-ossp";

create schema if not exists keeper;

create table if not exists keeper.usr(
    id uuid not null default uuid_generate_v4(),
    login varchar(64) not null,
    password varchar(128) not null,
    constraint usr_pkey primary key (id),
    constraint usr_login_uk unique (login)
);

create table if not exists keeper.client(
    id uuid not null default uuid_generate_v4(),
    user_id uuid not null,
    sync_tms timestamp not null default CURRENT_TIMESTAMP,
    constraint client_pkey primary key (id),
    constraint fk_client_usr_id foreign key(user_id) references keeper.usr(id)
);

create table if not exists keeper.cred(
    id uuid not null,
    login varchar(128) not null,
    password varchar(256) not null,
    user_id uuid not null,
    status varchar(8) not null default 'ACTIVE',
    modified_tms timestamp not null,
    constraint cred_pkey primary key (id)
);

create table if not exists keeper.txt(
    id uuid not null,
    val text not null,
    user_id uuid not null,
    status varchar(8) not null default 'ACTIVE',
    modified_tms timestamp not null,
    constraint txt_pkey primary key (id)
);

create table if not exists keeper.binary(
    id uuid not null,
    f_name varchar(1024) not null,
    "data" text not null,
    user_id uuid not null,
    status varchar(8) not null default 'ACTIVE',
    modified_tms timestamp not null,
    constraint binary_pkey primary key (id)
);

create table if not exists keeper.card(
    id uuid not null,
    num varchar(128) not null,
    cvc varchar(128) not null,
    holder_name varchar(128) not null,
    user_id uuid not null,
    status varchar(8) not null default 'ACTIVE',
    modified_tms timestamp not null,
    constraint card_pkey primary key (id)
);
-- +goose Down