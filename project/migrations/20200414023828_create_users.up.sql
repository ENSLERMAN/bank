create table bank.clients
(
    id         serial  not null,
    surname    varchar not null,
    name       varchar not null,
    patronymic varchar not null,
    login      varchar not null,
    password   varchar not null,
    passport   varchar not null
);

alter table bank.clients
    owner to ensler;

create unique index clients_id_uindex
    on bank.clients (id);

create unique index clients_login_uindex
    on bank.clients (login);

create unique index clients_passport_uindex
    on bank.clients (passport);

create unique index clients_password_uindex
    on bank.clients (password);
