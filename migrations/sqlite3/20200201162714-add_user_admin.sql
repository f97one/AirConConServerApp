
-- +migrate Up
/* app_user の digest_id を DROP して need_pw_change を ADD したいが、 */
/* SQLite は drop column をサポートしていないので、新しい構造でテーブルを作って */
/* 既存のレコードを複写したあと、古いテーブルを DROP して新しいテーブルを本の名前に */
/* rename to する方法をとる必要がある */
create table app_user_new (
    user_id integer primary key autoincrement,
    username text not null unique,
    password text not null default '',
    need_pw_change integer default 0
);
insert into app_user_new (user_id, username, password) select user_id, username, password from app_user;
drop table app_user;
alter table app_user_new rename to app_user;

/* ユーザー admin/admin を追加 */
insert into app_user (username, password, need_pw_change) values ('admin', '$2a$10$KyajCUdu7yo7XQKjh4eMtOyYPj1QFQvgmMkPY.KHAxztglP6oZuRe', 1);

-- +migrate Down
delete from app_user where username = 'admin';

create table app_user_old (
    user_id integer primary key autoincrement,
    username text not null unique,
    password text not null default '',
    digest_id text default ''
);
insert into app_user_old (user_id, username, password) select user_id, username, password from app_user;
drop table app_user;
alter table app_user_old rename to app_user;
