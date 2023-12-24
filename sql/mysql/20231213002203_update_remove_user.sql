-- +goose Up
-- +goose StatementBegin
create procedure if not exists DeleteUser (
  in p_user_oid varchar(36)
) begin
  delete from user as u
  where u.o_id=p_user_oid;
end;
-- +goose StatementEnd

-- +goose StatementBegin
create procedure UpdateUserRole (
  in p_user_oid varchar(36),
  in p_role varchar(30)
) begin
  update user u, user_role r
  set u.role=r.id
  where r.label=p_role and u.o_id=p_user_oid;
end;
-- +goose StatementEnd


-- +goose Down
drop procedure DeleteUser;
drop procedure UpdateUserRole;
