
-- +migrate Up
create table app_user (
    user_id integer primary key autoincrement,
    username text not null unique,
    password text not null unique,
    digest_id text default ''
);
create table jwt_token (
    user_id integer not null,
    generated_token text not null default '',
    expires_at text
);

-- +migrate Down
drop table jwt_token;
drop table app_user;