package gateways

import (
	"database/sql"
	"fmt"
	"strings"
	"user-management/internal/domain"
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

func (r ArticleRepository) Create(userOId string, article repo.ArticleData) error {
	fmt.Println("> repo: before create article query")
	query := fmt.Sprintf(`
		insert into article (o_id, user_id, theme, text)
		select "%v", id, "%v", "%v"
		from user
		where user.o_id="%v";`, article.OId, article.Theme, article.Text, userOId)
	if len(article.Tags) > 0 {
		query += fmt.Sprintf(`
		insert into tag (label)
		values %v as new
		on duplicate key update label=new.label;
		
		insert into article_tag (article_id, tag_id)
		select a.id as article_id, t.id as tag_id
		from article as a
		join tag t on a.o_id="%v" and t.label in (%v);`,
			tagsToString(article.Tags, "('", "')"),
			article.OId, tagsToString(article.Tags, "'", "'"))
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
		logrus.Infoln("> failed to scan article", err)
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
		on a.id=ats.article_id and ats.tag_id=t.id
		union
		select  a.id, a.o_id, a.user_id, a.theme, a.text, t.label as tags, a.created_at, a.updated_at
		from article a
		join (article_tag as ats join tag t)
		on a.id=ats.article_id  and ats.tag_id=t.id
		where a.id is null) as a,
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

	logrus.Infof("> scanned articles: %+v\n", res)
	return res, nil
}
