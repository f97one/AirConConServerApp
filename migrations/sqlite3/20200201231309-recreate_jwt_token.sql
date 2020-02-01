
-- +migrate Up
create table jwt_token_new (
    user_id integer primary key,
    generated_token text not null default '',
    expires_at text
);
drop table jwt_token;
alter table jwt_token_new rename to jwt_token;

-- +migrate Down
create table jwt_token_old (
    user_id integer not null,
    generated_token text not null default '',
    expires_at text
);
drop table jwt_token;
alter table jwt_token_old rename to jwt_token;