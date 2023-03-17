CREATE TABLE users (
    id serial not null unique,
    name text not null,
    username varchar(128) not null unique,
    password varchar(128) not null,
    email text not null unique
);