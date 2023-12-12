-- +goose Up

drop procedure IsOwnerOfArticle;
-- +goose StatementBegin
create procedure IsOwnerOfArticle (
  in article_oid varchar(36),
  in user_oid varchar(36)
) begin
  select a.o_id
  from user u
  inner join article a
  on u.id=a.user_id
  where a.o_id=article_oid and 
    u.o_id=user_oid
  group by a.o_id;
end
-- +goose StatementEnd
drop procedure CreateArticle;
-- +goose StatementBegin
create procedure if not exists CreateArticle (
  in p_user_oid varchar(36),
  in p_article_oid varchar(36),
  in p_theme varchar(200),
  in p_text varchar(200)
) begin
    insert into article (o_id, user_id, theme, text)
		select p_article_oid, u.id, p_theme, p_text
		from user as u
		where user.o_id=p_user_oid;

    select id
    from article
    where o_id=p_article_oid;
  end
-- +goose StatementEnd




-- +goose Down

drop procedure IsOwnerOfArticle;
-- +goose StatementBegin
create procedure IsOwnerOfArticle (
  in article_oid varchar(36),
  in username varchar(50)
) begin
  select a.o_id,  a.theme, a.text, group_concat(t.label) as tags, a.created_at, a.updated_at, a.status
  from user u
  inner join article a
  on u.id=a.user_id
  left join article_tag as ats 
  on a.id=ats.article_id
  left join tag t
  on ats.tag_id=t.id
  where a.o_id=article_oid and 
    u.username=username
  group by a.o_id;
end
-- +goose StatementEnd

drop procedure CreateArticle;
-- +goose StatementBegin
create procedure if not exists CreateArticle (
  in user_oid varchar(36),
  in article_oid varchar(36),
  in theme varchar(200),
  in text varchar(200)
) begin
    insert into article (o_id, user_id, theme, text)
		select article_oid, u.id, theme, text
		from user as u
		where user.o_id=user_oid;
  end
-- +goose StatementEnd