drop table if exists entries;
create table entries (
       id integer primary key autoincrement,
       title text not null,
       'text' text not null
);

drop table if exists users;
create table users (
       id integer primary key autoincrement,
       firstname text not null,
       lastname text not null,
       mail_address text,
       birthtdate timestamp,
       username text not null,
       password text not null
);
