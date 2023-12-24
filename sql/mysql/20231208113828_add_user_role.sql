-- +goose Up

alter table user
add column role int not null;

drop procedure CreateUser;
drop procedure GetUserByUsername;
drop procedure GetUserByOId;

-- +goose StatementBegin
create procedure if not exists CreateUser (
  in p_user_oid varchar(36),
  in p_username varchar(50),
  in p_pswd varchar(60),
  in p_role int
) begin
  insert into user (o_id, username, pswd, role)
	values (p_user_oid, p_username, p_pswd, p_role);
end
-- +goose StatementEnd

-- +goose StatementBegin
create procedure if not exists GetUserByUsername (
  in usernameV varchar(50)
) begin
  select o_id, username, pswd, role
  from user
  where user.username=usernameV;
end
-- +goose StatementEnd

-- +goose StatementBegin
create procedure if not exists GetUserByOId (
  in oidV varchar(36)
) begin
  select o_id, username, pswd, role
  from user
  where user.o_id=oidV;
end
-- +goose StatementEnd




-- +goose Down

alter table user
drop column role;

drop procedure CreateUser;
drop procedure GetUserByUsername;
drop procedure GetUserByOId;

-- +goose StatementBegin
create procedure if not exists CreateUser (
  in user_oid varchar(36),
  in usernameV varchar(50),
  in pswdV varchar(60)
) begin
  insert into user (o_id, username, pswd)
	values (user_oid, usernameV, pswdV);
end
-- +goose StatementEnd

-- +goose StatementBegin
create procedure if not exists GetUserByUsername (
  in usernameV varchar(50)
) begin
  select o_id, username, pswd
  from user
  where user.username=usernameV;
end
-- +goose StatementEnd

-- +goose StatementBegin
create procedure if not exists GetUserByOId (
  in oidV varchar(36)
) begin
  select o_id, username, pswd
  from user
  where user.o_id=oidV;
end
-- +goose StatementEnd
