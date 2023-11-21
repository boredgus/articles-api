create table user (
  id int auto_increment primary key,
  o_id varchar(36) unique not null,
  username varchar(50) unique not null check (username > ''),
  pswd varchar(60) not null check (pswd > '')
);
