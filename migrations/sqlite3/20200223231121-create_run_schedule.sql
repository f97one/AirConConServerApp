
-- +migrate Up
create table job_schedule (
    schedule_id text primary key not null,
    job_id integer not null,
    cmd_line text,
    run_at text
);

-- +migrate Down
drop table job_schedule;
