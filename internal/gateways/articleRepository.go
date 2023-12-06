package gateways

import (
	"database/sql"
	"strings"
	"time"
	"user-management/internal/domain"
	"user-management/internal/models"
	"user-management/internal/models/repo"
)

func NewArticleRepository(store Store) repo.ArticleRepository {
	return ArticleRepository{store: store}
}

type ArticleRepository struct {
	store Store
}

func tagsToArrayStr(tags []string) (res string) {
	for i, t := range tags {
		res += "'" + t + "'"
		if i != len(tags)-1 {
			res += ","
		}
	}
	return
}

func addTagsForArticleQuery(count int) (res string) {
	for i := 0; i < count; i++ {
		res += "call CreateTag(?);\n"
	}
	res += "call AddTagsToArticle(?,?);\n"
	return
}

func (r ArticleRepository) Create(userOId string, article repo.ArticleData) error {
	query := "call CreateArticle(?,?,?,?);\n"
	args := make([]any, 0, 4+len(article.Tags)+2)
	args = append(args, userOId, article.OId, article.Theme, article.Text)
	if len(article.Tags) > 0 {
		query += addTagsForArticleQuery(len(article.Tags))
		for i := range article.Tags {
			args = append(args, article.Tags[i])
		}
		args = append(args, article.OId, tagsToArrayStr(article.Tags))
	}
	rows, err := r.store.Query(query, args...)
	if err != nil {
		return err
	}
	rows.Close()
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

func (r ArticleRepository) Get(articleOId string) (a domain.Article, err error) {
	rows, err := r.store.Query(`call GetArticle(?);`, articleOId)
	if err != nil {
		return
	}
	if !rows.Next() {
		return a, models.ArticleNotFoundErr
	}
	a, err = r.scan(rows)
	if err != nil {
		return
	}
	rows.Close()
	return a, nil
}

func (r ArticleRepository) GetForUser(username string, page, limit int) ([]domain.Article, error) {
	rows, err := r.store.Query(`call GetArticlesForUser(?,?,?);`, username, page*limit, limit)
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
	rows.Close()
	return res, nil
}

func (r ArticleRepository) IsOwner(articleOId, username string) error {
	rows, err := r.store.Query(`call IsOwnerOfArticle(?,?);`, articleOId, username)
	if err != nil {
		return err
	}
	rows.Next()
	var res string
	err = rows.Scan(&res)
	if err != nil {
		return models.UserIsNotAnOwnerErr
	}
	rows.Close()
	return nil
}

func (r ArticleRepository) Update(article repo.ArticleData) (time.Time, error) {
	query := "call UpdateArticle(?,?,?);\n"
	args := make([]any, 0, 3+1+len(article.Tags)+2+1+1)
	args = append(args, article.Theme, article.Text, article.OId)

	if len(article.Tags) == 0 {
		query += "call RemoveAllTagsForArticle(?);\n"
		args = append(args, article.OId)
	} else {
		query += addTagsForArticleQuery(len(article.Tags)) + "call RemoveTagsForArticle(?,?);\n"
		for i := range article.Tags {
			args = append(args, article.Tags[i])
		}
		args = append(args, article.OId, tagsToArrayStr(article.Tags), article.OId, tagsToArrayStr(article.Tags))
	}
	query += "call GetTimeOfCreation(?);"
	rows, err := r.store.Query(query, append(args, article.OId))
	if err != nil {
		return time.Time{}, err
	}
	rows.Next()
	var timeOfCreation sql.NullTime
	err = rows.Scan(&timeOfCreation)
	if err != nil {
		return time.Time{}, err
	}
	rows.Close()
	return timeOfCreation.Time, nil
}
