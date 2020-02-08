
-- +migrate Up
create table schedule (
    schedule_id text primary key,
    name text not null,
    on_off integer,
    execute_time text default '00:00',
    script_id text not null
);
create table timing (
    schedule_id text primary key,
    weekday_id integer default 0
);

-- +migrate Down
drop table timing;
drop table schedule;