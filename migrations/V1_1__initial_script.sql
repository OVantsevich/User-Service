create table if not exists roles
(
    id      varchar(100)
        constraint Role_pk
            primary key,
    "name"  varchar(50)                               not null,
    deleted boolean                                   not null default false,
    created timestamp(6) default CURRENT_TIMESTAMP(6) not null,
    updated timestamp(6) default CURRENT_TIMESTAMP(6) not null
);

alter table roles
    owner to postgres;

create unique index if not exists roles_name_index
    on roles ("name");

insert into roles (id, "name")
values (gen_random_uuid(), 'admin');
insert into roles (id, "name")
values (gen_random_uuid(), 'user');

create table if not exists users
(
    id         varchar(100)
        constraint User_pk
            primary key,
    "login"    varchar(50)  not null,
    email      varchar(50)  not null,
    "password" varchar(200) not null,
    "name"     varchar(50)  not null,
    age        integer      not null,
    "role"     varchar(100) not null default 'user',
    token      varchar(200) not null default '',
    deleted    boolean      not null default false,
    created    timestamp(6)          default CURRENT_TIMESTAMP(6) not null,
    updated    timestamp(6)          default CURRENT_TIMESTAMP(6) not null
);

alter table users
    owner to postgres;

create unique index if not exists users_login_uindex
    on users ("login");