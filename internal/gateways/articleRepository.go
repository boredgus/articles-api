package gateways

import (
	"database/sql"
	"fmt"
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

func tagsToString(tags []string, leftWrapper, rightWrapper, separator string) (res string) {
	for i, t := range tags {
		res += leftWrapper + t + rightWrapper
		if i != len(tags)-1 {
			res += separator
		}
	}
	return
}

func addTagsForArticleQuery(articleOId string, tags []string) (res string) {
	for _, tag := range tags {
		res += fmt.Sprintf(`call CreateTag('%v');`, tag)
	}
	res += fmt.Sprintf(`call AddTagsToArticle("%v", "%v");`, articleOId, tagsToString(tags, "'", "'", ","))
	return
}

func (r ArticleRepository) Create(userOId string, article repo.ArticleData) error {
	query := fmt.Sprintf(`
		call CreateArticle('%v', '%v', '%v', '%v');`,
		userOId, article.OId, article.Theme, article.Text)
	if len(article.Tags) > 0 {
		query += addTagsForArticleQuery(article.OId, article.Tags)
	}
	rows, err := r.store.Query(query)
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
	query := fmt.Sprintf(`call UpdateArticle('%v','%v','%v');`,
		article.Theme, article.Text, article.OId)
	if len(article.Tags) == 0 {
		query += fmt.Sprintf(`call RemoveAllTagsForArticle('%v');`, article.OId)
	} else {
		query += addTagsForArticleQuery(article.OId, article.Tags) +
			fmt.Sprintf(`call RemoveTagsForArticle("%v","%v");`,
				article.OId, tagsToString(article.Tags, "'", "'", ","))
	}
	query += fmt.Sprintf(`call GetTimeOfCreation('%v');`, article.OId)
	rows, err := r.store.Query(query)
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
