create table if not exists rooms (
    id int unsigned not null auto_increment primary key,
    status int not null default 0,
    capacity int not null,
    host_id int unsigned not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp
);

create table if not exists room_players (
    id int unsigned not null auto_increment primary key,
    room_id int unsigned not null,
    player_id int unsigned not null,
    is_ready tinyint(1) not null default 0,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,
    foreign key (room_id) references rooms(id) on delete cascade
);