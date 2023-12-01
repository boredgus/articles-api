package gateways

import (
	"database/sql"
	"fmt"
	"strings"
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

func tagsToString(tags []string, leftWrapper, rightWrapper string) (res string) {
	for i, t := range tags {
		res += leftWrapper + t + rightWrapper
		if i != len(tags)-1 {
			res += ","
		}
	}
	return
}

func addTagsForArticle(OId string, tags []string) string {
	return fmt.Sprintf(`
		insert into tag (label)
		values %v as new
		on duplicate key update label=new.label;
		
		insert into article_tag (article_id, tag_id)
		select a.id as article_id, t.id as tag_id
		from article as a
		join tag t on a.o_id="%v" and t.label in (%v);`,
		tagsToString(tags, "('", "')"), OId, tagsToString(tags, "'", "'"))
}

func (r ArticleRepository) Create(userOId string, article repo.ArticleData) error {
	fmt.Println("> repo: before create article query")
	query := fmt.Sprintf(`
		insert into article (o_id, user_id, theme, text)
		select "%v", id, "%v", "%v"
		from user
		where user.o_id="%v";`, article.OId, article.Theme, article.Text, userOId)
	if len(article.Tags) > 0 {
		query += addTagsForArticle(article.OId, article.Tags)
	}
	_, err := r.store.Query(query)
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
	err := rows.Scan(&a.OId, &a.Theme, &a.Text, &tags, &a.CreatedAt, &updatedAt)
	if err != nil {
		return domain.Article{}, err
	}
	a.UpdatedAt = &updatedAt.Time
	if tags.Valid {
		a.Tags = strings.Split(tags.String, ",")
	} else {
		a.Tags = []string{}
	}
	return a, nil
}

func (r ArticleRepository) GetForUser(username string, page, limit int) ([]domain.Article, error) {
	rows, err := r.store.Query(`
	select a.o_id, a.theme, a.text, group_concat(a.tags), a.created_at, a.updated_at
	from
		(select a.id, a.o_id, a.user_id, a.theme, a.text, t.label as tags, a.created_at, a.updated_at
		from article a
		left join (article_tag as ats join tag t)
		on a.id=ats.article_id and ats.tag_id=t.id) as a,
		user as u
	where u.username=? and a.user_id=u.id
	group by a.id, a.o_id, a.theme, a.text, a.created_at, a.updated_at
	order by a.id desc
	limit ?, ?;`, username, page*limit, limit)
	if err != nil {
		fmt.Println("> failed to query articles", err)
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

func (r ArticleRepository) Update(article repo.ArticleData) error {
	/**
	update theme, text, status
	tags length == 0:
		remove all tags from article_tag
	tags length > 0:
		create tags
		remove tags from article_tag that are not equal to given
	*/
	query := fmt.Sprintf(`
		update article a
		set a.theme="%v" and a.text="%v" and a.status=1
		where a.o_id=%v;`, article.Theme, article.Text, article.OId)
	if len(article.Tags) == 0 {
		query += fmt.Sprintf(`
			delete ats
			from article_tag ats, article a
			where ats.article_id=a.id and a.o_id="%v";`, article.OId)
	} else {
		query += fmt.Sprintf(`%v
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
			addTagsForArticle(article.OId, article.Tags), article.OId, tagsToString(article.Tags, "'", "'"))
	}
	rows, err := r.store.Query(query)

	fmt.Println("> repo: update article err", err)
	cols, err := rows.Columns()
	fmt.Println("> repo: columns", cols, err)

	rows.Next()
	var res string
	fmt.Println("> repo: err invoke", rows.Err())
	err = rows.Scan(&res)

	fmt.Println("> repo: scan err and res", err, res)

	return nil
}
