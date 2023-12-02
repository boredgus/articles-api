package gateways

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"user-management/internal/domain"
	"user-management/internal/models"
	"user-management/internal/models/repo"

	"github.com/sirupsen/logrus"
)

func NewArticleRepository(store Store) repo.ArticleRepository {
	return ArticleRepository{store: store}
}

type ArticleRepository struct {
	store Store
}

func tagsToString(tags []string, leftWrapper, rightWrapper, separator string) (res string) {
	for i, t := range tags {
		res += leftWrapper + t + rightWrapper
		if i != len(tags)-1 {
			res += separator
		}
	}
	return
}

func addTagsForArticle(articleOId string, tags []string) string {
	return fmt.Sprintf(`
	-- create new tags without duplication
	insert into tag (label)
	select * 
	from (%v) as tmp
	where not exists (
		select label
		from tag
		where label=tmp.label
	);

	-- bind tags to article without duplication of existed bindings
	insert into article_tag (article_id, tag_id)
	select new_.article_id, new_.tag_id
	from 
		(select a.id article_id, t.id tag_id, t.label
		from  article a, tag t
		where a.o_id="%v" and 
			t.label in (%v)) as new_
	where new_.label not in (
		select t.label
		from  article a, article_tag ats
		join tag t
		on ats.tag_id=t.id
		where a.o_id="%v" and 
			ats.article_id=a.id);`,
		tagsToString(tags, "select '", "' label", " union "), articleOId,
		tagsToString(tags, "'", "'", ","), articleOId,
	)
}

func (r ArticleRepository) Create(userOId string, article repo.ArticleData) error {
	query := fmt.Sprintf(`
		insert into article (o_id, user_id, theme, text)
		select "%v", id, "%v", "%v"
		from user
		where user.o_id="%v";`, article.OId, article.Theme, article.Text, userOId)
	if len(article.Tags) > 0 {
		query += addTagsForArticle(article.OId, article.Tags)
	}
	rows, err := r.store.Query(query)
	rows.Close()
	if err != nil {
		logrus.Infoln("failed to create article", err)
		return err
	}
	return nil
}
func (r ArticleRepository) scan(rows *sql.Rows) (domain.Article, error) {
	var a domain.Article
	var tags sql.NullString
	var updatedAt sql.NullTime
	err := rows.Scan(&a.OId, &a.Theme, &a.Text, &tags, &a.CreatedAt, &updatedAt, &a.Status)
	if err != nil {
		return domain.Article{}, err
	}
	if updatedAt.Valid {
		a.UpdatedAt = &updatedAt.Time
	}
	if tags.Valid {
		a.Tags = strings.Split(tags.String, ",")
	} else {
		a.Tags = []string{}
	}
	return a, nil
}

func (r ArticleRepository) GetForUser(username string, page, limit int) ([]domain.Article, error) {
	rows, err := r.store.Query(`
	select a.o_id, a.theme, a.text, group_concat(a.tags), a.created_at, a.updated_at, a.status
	from
		(select a.id, a.o_id, a.user_id, a.theme, a.text, t.label as tags, a.created_at, a.updated_at, a.status
		from article a
		left join (article_tag as ats join tag t)
		on a.id=ats.article_id and ats.tag_id=t.id) as a,
		user as u
	where u.username=? and a.user_id=u.id
	group by a.id, a.o_id, a.theme, a.text, a.created_at, a.updated_at, a.status
	order by a.id desc
	limit ?, ?;`, username, page*limit, limit)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	res := make([]domain.Article, 0, limit)
	for rows.Next() {
		article, err := r.scan(rows)
		if err != nil {
			return []domain.Article{}, err
		}
		res = append(res, article)
	}

	return res, nil
}

func (r ArticleRepository) IsOwner(articleOId, username string) error {
	rows, err := r.store.Query(`
	select a.o_id
	from article a, user u
	where u.username=? and
		u.id=a.user_id and a.o_id=?;`, username, articleOId)
	defer rows.Close()
	if err != nil {
		return err
	}
	rows.Next()
	var res string
	err = rows.Scan(&res)
	if err != nil {
		return models.UserIsNotAnOwnerErr
	}
	return nil
}

func (r ArticleRepository) Update(article repo.ArticleData) (time.Time, error) {
	query := fmt.Sprintf(`
		-- update article
		update article a
		set a.theme="%v", a.text="%v", a.status=1
		where a.o_id="%v";`, article.Theme, article.Text, article.OId)
	if len(article.Tags) == 0 {
		query += fmt.Sprintf(`
			-- delete all tags for this article
			delete ats
			from article_tag ats, article a
			where ats.article_id=a.id and a.o_id="%v";`, article.OId)
	} else {
		query += addTagsForArticle(article.OId, article.Tags) +
			fmt.Sprintf(`
				-- delete tags that are not in given data
				delete ats
				from
					article_tag as ats,
					(select a.id as article_id, t.id as tag_id, t.label as tag
					from article a
					left join (article_tag as ats join tag t)
					on a.id=ats.article_id and ats.tag_id=t.id
						where a.o_id="%v") as tags
				where ats.article_id=tags.article_id and
					tags.tag not in (%v);`,
				article.OId, tagsToString(article.Tags, "'", "'", ","))
	}
	query += fmt.Sprintf(`
	-- get time of creation
	select created_at
	from article
	where o_id="%v";
	`, article.OId)
	rows, err := r.store.Query(query)
	defer rows.Close()
	if err != nil {
		return time.Time{}, err
	}
	rows.Next()
	var timeOfCreation sql.NullTime
	err = rows.Scan(&timeOfCreation)
	if err != nil {
		return time.Time{}, err
	}
	return timeOfCreation.Time, nil
}
