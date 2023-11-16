CREATE TABLE user (
  id int auto_increment,
  o_id varchar(36) unique not null,
  username varchar(100) unique not null,
  pswd varchar(200) not null,
  primary key (id)
);

drop table user_profile;
drop table user;

create table user (
	id int auto_increment,
    username varchar(255) unique not null,
	pswd VARCHAR(255) not null,
	primary key (id)
);

-- create table user_profile (
-- 	id int,
-- 	o_id varchar(36) not null unique,
-- 	nickname varchar(255) not null unique,
--     first_name varchar(255) default "",
--     last_name varchar(255) default "",
--     created_at timestamp default current_timestamp,
--     updated_at timestamp,
--     deleted_at timestamp,
--     primary key (id),
--     foreign key (id) references user (id) on delete cascade
-- );
