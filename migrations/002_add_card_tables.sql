create table if not exists cards (
    id int unsigned not null auto_increment primary key,
    name varchar(50) not null,
    type tinyint not null,
    cost int not null,
    description varchar(500) not null
);

create table if not exists player_card_states (
    id int unsigned not null auto_increment primary key,
    game_id int unsigned not null,
    player_id int unsigned not null,
    stage tinyint not null default 0,
    hand_card_ids json,
    deck_card_ids json,
    discard_card_ids json,
    hand_basic_count int not null default 0,
    hand_skill_count int not null default 0,
    deck_basic_count int not null default 0,
    deck_skill_count int not null default 0,
    basic_card_played tinyint(1) not null default 0
);