create table article (
  id int auto_increment primary key,
  o_id varchar(36) unique not null,
  user_id int,
  theme varchar(200) not null check (theme > ''),
  text varchar(500) default "",
  created_at timestamp default current_timestamp,
  updated_at timestamp on update current_timestamp,
  status int not null default 0,
  foreign key (user_id) references user(id) on delete cascade
);

create table tag (
  id int auto_increment primary key,
  label varchar(100) not null unique check (label > '')
);

create table article_tag (
  article_id int,
  tag_id int,
  foreign key (article_id) references article(id) on delete cascade,
  foreign key (tag_id) references tag(id) on delete cascade
);
