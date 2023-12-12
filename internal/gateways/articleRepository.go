package gateways

import (
	"database/sql"
	"strings"
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

func (r ArticleRepository) CreateArticle(userOId string, article repo.ArticleData) error {
	rows, err := r.store.Query("call CreateArticle(?,?,?,?);", userOId, article.OId, article.Theme, article.Text)
	if err != nil {
		return err
	}
	rows.Close()
	if len(article.Tags) > 0 {
		return r.AddTagsForArticle(article.OId, article.Tags)
	}
	return nil
}
func (r ArticleRepository) DeleteArticle(oid string, tags []string) error {
	rows, err := r.store.Query("call DeleteArticle(?);", oid)
	if err != nil {
		return err
	}
	rows.Close()
	if len(tags) > 0 {
		return r.RemoveTagsFromArticle(oid, tags)
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

func (r ArticleRepository) IsOwner(articleOId, userOId string) error {
	rows, err := r.store.Query(`call IsOwnerOfArticle(?,?);`, articleOId, userOId)
	if err != nil {
		return err
	}
	if !rows.Next() {
		return models.ArticleNotFoundErr
	}
	rows.Close()
	return nil
}
func (r ArticleRepository) UpdateArticle(oid, theme, text string) error {
	rows, err := r.store.Query("call UpdateArticle(?,?,?);", oid, theme, text)
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}
func (r ArticleRepository) AddTagsForArticle(articleOId string, tags []string) error {
	var query string
	args := make([]any, 0, len(tags)+2)
	for _, t := range tags {
		query += "call CreateTag(?);\n"
		args = append(args, t)
	}
	query += "call AddTagsToArticle(?,?);\n"
	args = append(args, articleOId, tagsToArrayStr(tags))
	rows, err := r.store.Query(query, args...)
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}
func (r ArticleRepository) RemoveTagsFromArticle(articleOId string, tags []string) error {
	rows, err := r.store.Query("call RemoveTagsForArticle(?,?);", articleOId, tagsToArrayStr(tags))
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}
