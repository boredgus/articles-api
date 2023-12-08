-- +goose Up
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

-- +goose Down
drop procedure CreateUser;
drop procedure GetUserByUsername;
drop procedure GetUserByOId;
