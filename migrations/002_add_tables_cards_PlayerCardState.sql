create table if not exists cards (
    id int unsigned not null auto_increment primary key,
    name varchar(50) not null,
    type tinyint not null,
    cost int not null,
    description varchar(500) not null,
);