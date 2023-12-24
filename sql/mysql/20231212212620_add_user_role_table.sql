-- +goose Up
create table user_role (
  id int primary key,
  label varchar(30) unique not null check (label > '')
);

insert into user_role (id,label)
values (0,"user"), (1,"moderator"), (2,"admin");

alter table user
add constraint user_role_fk
foreign key (role) references user_role (id) on delete cascade;

drop procedure CreateUser;
-- +goose StatementBegin
create procedure CreateUser (
  in p_user_oid varchar(36),
  in p_username varchar(50),
  in p_pswd varchar(60),
  in p_role varchar(30)
) begin
  insert into user (o_id, username, pswd, role)
  select p_user_oid, p_username, p_pswd, r.id
  from user_role as r
  where r.label=p_role;
end
-- +goose StatementEnd

drop procedure GetUserByUsername;
-- +goose StatementBegin
create procedure if not exists GetUserByUsername (
  in p_username varchar(50)
) begin
  select u.o_id, u.username, u.pswd, r.label
  from user as u
  inner join user_role as r
  on u.role=r.id
  where u.username=p_username;
end
-- +goose StatementEnd


-- +goose Down
alter table user
drop constraint user_role_fk;

drop table user_role;

drop procedure CreateUser;
-- +goose StatementBegin
create procedure CreateUser (
  in p_user_oid varchar(36),
  in p_username varchar(50),
  in p_pswd varchar(60),
  in p_role int
) begin
  insert into user (o_id, username, pswd, role)
	values (p_user_oid, p_username, p_pswd, p_role);
end
-- +goose StatementEnd

drop procedure GetUserByUsername;
-- +goose StatementBegin
create procedure if not exists GetUserByUsername (
  in usernameV varchar(50)
) begin
  select o_id, username, pswd, role
  from user
  where user.username=usernameV;
end
-- +goose StatementEnd
