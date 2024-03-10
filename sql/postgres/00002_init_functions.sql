-- +goose Up

-- +goose StatementBegin
CREATE OR REPLACE PROCEDURE articlesdb.CreateArticle (
  p_user_oid UUID,
  p_article_oid UUID,
  p_theme varchar(200),
  p_text varchar(200)
)
LANGUAGE plpgsql
AS $$
BEGIN
  insert into articlesdb."article" (o_id, user_id, theme, text)
  select p_article_oid, u.id, p_theme, p_text
  from articlesdb."user" as u
  where u.o_id=p_user_oid;
END $$
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE PROCEDURE articlesdb.CreateTag (
  p_label varchar(100)
)
LANGUAGE plpgsql
AS $$
BEGIN
	if not exists (select label from articlesdb."tag" where label=p_label)
	then 
		insert into tag (label) values (p_label);
	end if;
end $$
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE PROCEDURE articlesdb.CreateUser (
  p_user_oid UUID,
  p_username varchar(50),
  p_pswd varchar(60),
  p_role varchar(30)
)
LANGUAGE plpgsql
AS $$
BEGIN
  insert into articlesdb."user" (o_id, username, pswd, role)
  select p_user_oid, p_username, p_pswd, r.id
  from articlesdb."user_role" as r
  where r.label=p_role;
end $$
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE PROCEDURE articlesdb.DeleteArticle (
  p_article_oid UUID
)
LANGUAGE plpgsql
AS $$
BEGIN
  update articlesdb."article"
  set status=-1
  where o_id=p_article_oid;
END $$
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE PROCEDURE articlesdb.DeleteUser (
  p_user_oid UUID
)
LANGUAGE plpgsql
AS $$
BEGIN
  delete from articlesdb."user" as u
  where u.o_id=p_user_oid;
END $$
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TYPE articlesdb.article_data AS (
	o_id UUID,
	theme varchar(200),
	text varchar(500),
	tags text,
	created_at timestamp with time zone,
	updated_at timestamp with time zone,
	status integer
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE FUNCTION articlesdb.GetArticle (
  p_article_oid UUID
)
RETURNS SETOF articlesdb.article_data
LANGUAGE plpgsql
AS $$
BEGIN
  RETURN QUERY
    select a.o_id,  a.theme, a.text, string_agg(t.label,',') as tags, a.created_at, a.updated_at, a.status
    from articlesdb."article" as a
    left join articlesdb."article_tag" as ats 
    on a.id=ats.article_id
    left join articlesdb."tag" as t
    on ats.tag_id=t.id
    where a.o_id=p_article_oid and a.status != -1
    group by a.o_id,a.theme, a.text, a.created_at, a.updated_at, a.status;
END $$
-- +goose StatementEnd


-- +goose StatementBegin
CREATE FUNCTION articlesdb.GetArticlesForUser (
  p_username varchar(50),
  p_page integer,
  p_limit integer
)
RETURNS SETOF articlesdb.article_data
LANGUAGE plpgsql
AS $$
BEGIN
  RETURN QUERY 
	  select a.o_id, a.theme, a.text, string_agg(t.label,',') as tags, a.created_at, a.updated_at, a.status
	  from articlesdb."user" u
	  inner join articlesdb."article" a
	  on u.id=a.user_id
	  left join articlesdb."article_tag" as ats 
	  on a.id=ats.article_id
	  left join articlesdb."tag" t 
	  on ats.tag_id=t.id
	  where u.username=p_username and a.status != -1
	  group by a.id, a.o_id, a.theme, a.text, a.created_at, a.updated_at, a.status
	  order by a.id desc
	  limit p_limit offset p_page * p_limit;
END $$
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TYPE articlesdb.user_data AS (
	o_id UUID,
	username varchar(50),
  pswd varchar(60),
  role varchar(30)
);
-- +goose StatementEnd


-- +goose StatementBegin
CREATE FUNCTION articlesdb.GetUserByOId (
  p_user_oid UUID
)
RETURNS SETOF articlesdb.user_data
LANGUAGE plpgsql
AS $$
BEGIN
  RETURN QUERY
    select u.o_id, u.username, u.pswd, r.label
    from articlesdb."user" as u
    inner join articlesdb."user_role" as r
    on u.role=r.id
    where u.o_id=p_user_oid;
END $$
-- +goose StatementEnd

-- +goose StatementBegin
CREATE FUNCTION articlesdb.GetUserByUsername (
  p_username varchar(50)
)
RETURNS SETOF articlesdb.user_data
LANGUAGE plpgsql
AS $$
BEGIN
  RETURN QUERY
    select u.o_id, u.username, u.pswd, r.label
    from articlesdb."user" as u
    inner join articlesdb."user_role" as r
    on u.role=r.id
    where u.username=p_username;
END $$
-- +goose StatementEnd

-- +goose StatementBegin
CREATE FUNCTION articlesdb.IsOwnerOfArticle (
  p_article_oid UUID,
  p_user_oid UUID
)
RETURNS SETOF UUID
LANGUAGE plpgsql
AS $$
BEGIN
  RETURN QUERY
    select a.o_id
    from articlesdb."user" as u
    inner join articlesdb."article" as a
    on u.id=a.user_id
    where a.o_id=p_article_oid and 
      u.o_id=p_user_oid and a.status != -1
    group by a.o_id;
END $$
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE PROCEDURE articlesdb.UpdateArticle (
  p_article_oid UUID,
  p_theme varchar(200),
  p_text varchar(200)
)
LANGUAGE plpgsql
AS $$
BEGIN
  update articlesdb."article"
  set theme=p_theme, text=p_text, status=1
  where o_id=p_article_oid;
END $$
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE PROCEDURE articlesdb.UpdateUserRole (
  p_user_oid UUID,
  p_role varchar(30)
)
LANGUAGE plpgsql
AS $$
BEGIN
  update articlesdb."user"
  set role=articlesdb.user_role.id
  from articlesdb."user_role" 
  where articlesdb.user_role.label=p_role and articlesdb.user.o_id=p_user_oid;
END $$
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE PROCEDURE articlesdb.AddTagsToArticle (
  p_article_oid UUID,
  p_tag_array text
)
LANGUAGE plpgsql
AS $$
begin
	execute CONCAT('insert into articlesdb."article_tag" (article_id, tag_id)
    select a.id, t.id
    from articlesdb."tag" as t, articlesdb."article" as a
    where a.o_id=''',p_article_oid,''' and 
      t.label in (',p_tag_array,');');
END $$
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE PROCEDURE articlesdb.RemoveTagsForArticle (
  p_article_oid UUID,
  p_tag_array text
)
LANGUAGE plpgsql
AS $$
BEGIN
  execute CONCAT('delete from articlesdb."article_tag" as ats
    using articlesdb."article" as a,
      articlesdb."tag" as t
    where a.id=ats.article_id and ats.tag_id=t.id and
      a.o_id=''',p_article_oid,''' and t.label in (',p_tag_array,');');
END $$
-- +goose StatementEnd



-- +goose Down

DROP PROCEDURE articlesdb.CreateArticle;
DROP PROCEDURE articlesdb.CreateTag;
DROP PROCEDURE articlesdb.CreateUser;
DROP PROCEDURE articlesdb.DeleteArticle;
DROP PROCEDURE articlesdb.DeleteUser;
DROP FUNCTION articlesdb.GetArticle;
DROP FUNCTION articlesdb.GetArticlesForUser;
DROP FUNCTION articlesdb.GetUserByOId;
DROP FUNCTION articlesdb.GetUserByUsername;
DROP FUNCTION articlesdb.IsOwnerOfArticle;
DROP PROCEDURE articlesdb.UpdateArticle;
DROP PROCEDURE articlesdb.UpdateUserRole;
DROP PROCEDURE articlesdb.AddTagsToArticle;
DROP PROCEDURE articlesdb.RemoveTagsForArticle;
DROP TYPE articlesdb.article_data;
DROP TYPE articlesdb.user_data;
