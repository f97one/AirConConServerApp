
-- +migrate Up
create table scripts (
    script_id text primary key,
    gpio integer default 0,
    script_name text not null unique default '',
    freq real default 38
);

-- +migrate Down
drop table scripts;