
-- +migrate Up
alter table app_user add column admin_flag integer default 0;
update app_user set admin_flag = 1 where username = 'admin';

-- +migrate Down
create table app_user_old (
    user_id integer primary key autoincrement,
    username text not null unique,
    password text not null default '',
    need_pw_change integer default 0
);
insert into app_user_old (user_id, username, password, need_pw_change)
select user_id, username, password, need_pw_change from app_user;
drop table app_user;
alter table app_user_old rename to app_user;