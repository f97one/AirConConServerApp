
-- +migrate Up
create table timing_new (
    schedule_id text not null,
    weekday_id integer default 0,
    primary key (schedule_id, weekday_id)
);
insert into timing_new (schedule_id, weekday_id) select schedule_id, weekday_id from timing;
drop table timing;
alter table timing_new rename to timing;

-- +migrate Down
create table timing_old (
    schedule_id text primary key,
    weekday_id integer default 0
);
insert into timing_old (schedule_id, weekday_id) select schedule_id, weekday_id from timing;
drop table timing;
alter table timing_old rename to timing;