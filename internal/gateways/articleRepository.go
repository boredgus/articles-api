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

func (r ArticleRepository) IsOwner(articleOId, username string) (domain.Article, error) {
	rows, err := r.store.Query(`call IsOwnerOfArticle(?,?);`, articleOId, username)
	if err != nil {
		return domain.Article{}, err
	}
	rows.Next()
	article, err := r.scan(rows)
	if err != nil {
		return domain.Article{}, models.UserIsNotAnOwnerErr
	}
	rows.Close()
	return article, nil
}

func (r ArticleRepository) Update(newA repo.ArticleData, oldA repo.ArticleData) error {
	query := "call UpdateArticle(?,?,?);\n"
	tagsToRemove, tagsToAdd := oldA.CompareTags(newA.Tags)
	args := make([]any, 0, 20)
	args = append(args, newA.OId, newA.Theme, newA.Text)
	if len(tagsToRemove) > 0 {
		query += "call RemoveTagsForArticle(?,?);\n"
		args = append(args, newA.OId, tagsToArrayStr(tagsToRemove))
	}
	if len(tagsToAdd) > 0 {
		query += addTagsForArticleQuery(len(tagsToAdd))
		for _, tag := range tagsToAdd {
			args = append(args, tag)
		}
		args = append(args, newA.OId, tagsToArrayStr(tagsToAdd))
	}
	rows, err := r.store.Query(query, args...)
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}
