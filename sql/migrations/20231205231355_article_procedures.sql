-- +goose Up
-- +goose StatementBegin
create procedure if not exists CreateTag (in tag varchar(100))
begin
  insert into tag (label)
  select tmp.label
  from (select tag label) as tmp
  where not exists (
    select label
    from tag
    where label=tmp.label
  );
end
-- +goose StatementEnd

-- +goose StatementBegin
create procedure if not exists AddTagsToArticle (
  in article_oid varchar(36),
  in tagArray varchar(500)
) begin
    set @query = concat('insert into article_tag (article_id, tag_id)
    select a.id, t.id
    from tag t, article a
    where a.o_id="',article_oid, '" and 
      t.label in (',tagArray,');');
    prepare stmt from @query;
    execute stmt;
    deallocate prepare stmt;
  end
-- +goose StatementEnd

-- +goose StatementBegin
create procedure if not exists CreateArticle (
  in user_oid varchar(36),
  in article_oid varchar(36),
  in theme varchar(200),
  in text varchar(200)
) begin
    insert into article (o_id, user_id, theme, text)
		select article_oid, id, theme, text
		from user
		where user.o_id=user_oid;
  end
-- +goose StatementEnd

-- +goose StatementBegin
create procedure if not exists GetArticle (
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

  
-- +goose StatementBegin
create procedure if not exists IsOwnerOfArticle (
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

-- +goose StatementBegin
create procedure if not exists UpdateArticle (
  in article_oid varchar(36),
  in theme varchar(200),
  in text varchar(200)
) begin
    update article a
		set a.theme=theme, a.text=text, a.status=1
		where a.o_id=article_oid;
  end
-- +goose StatementEnd

-- +goose StatementBegin
create procedure if not exists RemoveTagsForArticle (
  in article_oid varchar(36),
  in tagArray varchar(500)
) begin
    SET @query = CONCAT('delete ats
      from article a
      inner join article_tag ats
      on a.id=ats.article_id
      inner join tag t
      on ats.tag_id=t.id 
      where a.o_id="',article_oid,'" and t.label in ("',tagArray,'");');
    prepare stmt from @query;
    execute stmt;
    deallocate prepare stmt;
  end
-- +goose StatementEnd

-- +goose Down
drop procedure CreateTag;
drop procedure AddTagsToArticle;
drop procedure CreateArticle;
drop procedure GetArticle;
drop procedure GetArticlesForUser;
drop procedure IsOwnerOfArticle;
drop procedure RemoveTagsForArticle;
drop procedure UpdateArticle;

