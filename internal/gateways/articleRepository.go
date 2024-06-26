package gateways

import (
	"a-article/internal/domain"
	"a-article/internal/models"
	"a-article/internal/models/repo"
	"database/sql"
	"fmt"
	"strings"
)

func NewArticleRepository(mainStore, statsStore Store) repo.ArticleRepository {
	return &ArticleRepository{main: mainStore, stats: statsStore}
}

type ArticleRepository struct {
	main  Store
	stats Store
}

func arrayToStr(arr []string, isDouble bool) string {
	strBuilder := strings.Builder{}
	for i, t := range arr {
		if i > 0 {
			strBuilder.WriteString(",")
		}
		if isDouble {
			strBuilder.WriteString("''" + t + "''")
		} else {
			strBuilder.WriteString("'" + t + "'")
		}
	}
	return strBuilder.String()
}

func (r *ArticleRepository) CreateArticle(userOId string, article repo.ArticleData) error {
	rows, err := r.main.Query("call articlesdb.CreateArticle($1,$2,$3,$4);", userOId, article.OId, article.Theme, article.Text)
	if err != nil {
		return err
	}
	rows.Close()
	if len(article.Tags) > 0 {
		return r.AddTagsForArticle(article.OId, article.Tags)
	}
	return nil
}
func (r *ArticleRepository) DeleteArticle(oid string, tags []string) error {
	rows, err := r.main.Query("call articlesdb.DeleteArticle($1);", oid)
	if err != nil {
		return err
	}
	rows.Close()
	if len(tags) > 0 {
		return r.RemoveTagsFromArticle(oid, tags)
	}
	return nil
}
func (r *ArticleRepository) scan(rows Rows) (domain.Article, error) {
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

func (r *ArticleRepository) GetArticle(articleOId string) (a domain.Article, err error) {
	rows, err := r.main.Query(`select * from articlesdb.GetArticle($1);`, articleOId)
	if err != nil {
		return
	}
	if !rows.Next() {
		rows.Close()
		return a, models.NotFoundErr
	}
	a, err = r.scan(rows)
	if err != nil {
		return
	}
	rows.Close()
	return a, nil
}

func (r *ArticleRepository) GetForUser(username string, page, limit int) ([]domain.Article, error) {
	rows, err := r.main.Query(`select * from articlesdb.GetArticlesForUser($1,$2,$3);`, username, page*limit, limit)
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

func (r *ArticleRepository) IsOwner(articleOId, userOId string) error {
	rows, err := r.main.Query(`select * from articlesdb.IsOwnerOfArticle($1,$2);`, articleOId, userOId)
	if err != nil {
		return err
	}
	if !rows.Next() {
		return models.NotFoundErr
	}
	rows.Close()
	return nil
}
func (r *ArticleRepository) UpdateArticle(oid, theme, text string) error {
	rows, err := r.main.Query("call articlesdb.UpdateArticle($1,$2,$3);", oid, theme, text)
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}
func (r *ArticleRepository) AddTagsForArticle(articleOId string, tags []string) error {
	var builder strings.Builder
	for _, t := range tags {
		builder.WriteString(fmt.Sprintf("call articlesdb.CreateTag('%v');\n", t))
	}
	builder.WriteString(fmt.Sprintf("call articlesdb.AddTagsToArticle('%v','%v');\n", articleOId, arrayToStr(tags, true)))
	rows, err := r.main.Query(builder.String())
	if err != nil {
		return err
	}
	return rows.Close()
}
func (r *ArticleRepository) RemoveTagsFromArticle(articleOId string, tags []string) error {
	rows, err := r.main.Query("call articlesdb.RemoveTagsForArticle($1,$2);", articleOId, arrayToStr(tags, false))
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}

func (r *ArticleRepository) GetReactionsFor(articleOIds ...string) (repo.ArticleReactions, error) {
	rows, err := r.stats.Query(`SELECT article_id, reaction, sum(votes)
		FROM article_reaction FINAL
		WHERE article_id IN (?)
		GROUP BY article_id,reaction;`, articleOIds)
	if err != nil {
		return nil, err
	}
	reactions := repo.ArticleReactions{}
	for rows.Next() {
		var articleOId, reaction string
		var count int
		err := rows.Scan(&articleOId, &reaction, &count)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			continue
		}
		if reactions[articleOId] == nil {
			reactions[articleOId] = domain.ArticleReactions{}
		}
		reactions[articleOId][reaction] = int32(count)
	}
	rows.Close()
	return reactions, nil
}
func (r *ArticleRepository) GetCurrentReaction(raterOId, articleOId string) (string, error) {
	rows, err := r.stats.Query(`
		SELECT reaction, votes
		FROM article_reaction FINAL
		WHERE (article_id = ?) AND (rater_id = ?)`, articleOId, raterOId)
	if err != nil {
		return "", err
	}
	if !rows.Next() {
		rows.Close()
		return "", models.NotFoundErr
	}
	var reaction string
	var n int
	if err = rows.Scan(&reaction, &n); err != nil {
		return "", err
	}
	return reaction, rows.Close()
}
func (r *ArticleRepository) UpdateReaction(raterOId, articleOId, reaction string, count int) error {
	rows, err := r.stats.Query(`
		INSERT INTO article_reaction (article_id,rater_id,reaction,votes)
		VALUES (?,?,?,?)`, articleOId, raterOId, reaction, count)
	if err != nil {
		return err
	}
	return rows.Close()
}
