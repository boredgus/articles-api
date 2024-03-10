-- +goose Up
-- +goose StatementBegin
insert into articlesdb."user_role" (id,label)
values (0,'user'),(1,'moderator'),(2,'admin');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete from articlesdb."user_role";
-- +goose StatementEnd
