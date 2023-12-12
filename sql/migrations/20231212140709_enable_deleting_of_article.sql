-- +goose Up
-- +goose StatementBegin
create procedure DeleteArticle (
  in p_article_oid varchar(36)
) begin
  update article a
  set a.theme=theme, a.text=text, a.updated_at=CURRENT_TIMESTAMP(), a.status=-1
  where a.o_id=p_article_oid;
end
-- +goose StatementEnd

drop procedure GetArticle;
-- +goose StatementBegin
create procedure GetArticle (
  in p_article_oid varchar(36)
) begin
  select a.o_id,  a.theme, a.text, group_concat(t.label) as tags, a.created_at, a.updated_at, a.status
  from article a
  left join article_tag as ats 
  on a.id=ats.article_id
  left join tag t
  on ats.tag_id=t.id
  where a.o_id=p_article_oid and a.status!=-1
  group by a.o_id;
end
-- +goose StatementEnd

drop procedure GetArticlesForUser;
-- +goose StatementBegin
create procedure if not exists GetArticlesForUser (
  in p_username varchar(50),
  in p_page int,
  in p_limit int 
) begin
    select a.o_id, a.theme, a.text, group_concat(t.label) as tags, a.created_at, a.updated_at, a.status
    from user u
    inner join article a
    on u.id=a.user_id
    left join article_tag as ats 
    on a.id=ats.article_id
    left join tag t 
    on ats.tag_id=t.id
    where u.username=p_username and a.status!=-1
    group by a.o_id
    order by a.id desc
    limit p_page,p_limit;
  end
-- +goose StatementEnd

drop procedure IsOwnerOfArticle;
-- +goose StatementBegin
create procedure IsOwnerOfArticle (
  in p_article_oid varchar(36),
  in p_user_oid varchar(36)
) begin
  select a.o_id
  from user u
  inner join article a
  on u.id=a.user_id
  where a.o_id=p_article_oid and 
    u.o_id=p_user_oid and a.status!=-1
  group by a.o_id;
end
-- +goose StatementEnd


-- +goose Down
drop procedure DeleteArticle;
drop procedure GetArticle;
-- +goose StatementBegin
create procedure GetArticle (
  in article_oid varchar(36)
) begin
  select a.o_id,  a.theme, a.text, group_concat(t.label) as tags, a.created_at, a.updated_at, a.status
  from article a
  left join article_tag as ats 
  on a.id=ats.article_id
  left join tag t
  on ats.tag_id=t.id
  where a.o_id=article_oid
  group by a.o_id;
end
-- +goose StatementEnd

drop procedure GetArticlesForUser;
-- +goose StatementBegin
create procedure if not exists GetArticlesForUser (
  in username varchar(50),
  in pageV int,
  in limitV int 
) begin
    select a.o_id, a.theme, a.text, group_concat(t.label) as tags, a.created_at, a.updated_at, a.status
    from user u
    inner join article a
    on u.id=a.user_id
    left join article_tag as ats 
    on a.id=ats.article_id
    left join tag t 
    on ats.tag_id=t.id
    where u.username=username
    group by a.o_id
    order by a.id desc
    limit pageV,limitV;
  end
-- +goose StatementEnd

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