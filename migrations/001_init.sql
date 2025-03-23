create database moonlight;

use moonlight;

create table if not exists users (
  id bigint unsigned not null auto_increment primary key,
  username varchar(255) not null unique,
  password varchar(255) not null,
  role varchar(50) not null default 'player',
  created_at timestamp not null default current_timestamp,
  updated_at timestamp not null default current_timestamp on update current_timestamp
);