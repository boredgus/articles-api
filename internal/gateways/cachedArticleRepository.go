package gateways

import (
	"fmt"
	"user-management/internal/domain"
	"user-management/internal/models/repo"

	"github.com/sirupsen/logrus"
)

func NewCachedArticleRepository(repo repo.ArticleRepository, cacheStore CacheStore) repo.ArticleRepository {
	return &CachedArticleRepository{
		ArticleRepository: repo,
		cache:             cacheStore,
		serializer:        NewJSONSerializer[domain.Article](),
	}
}

func articleKey(oid string) string {
	return fmt.Sprintf("article:%v", oid)
}
func articlesKey(username string, page, limit int) string {
	return fmt.Sprintf("articles:%v:%v:%v", username, page, limit)
}

type CachedArticleRepository struct {
	repo.ArticleRepository
	cache      CacheStore
	serializer Serializer[domain.Article]
}

func (r *CachedArticleRepository) GetArticle(articleOId string) (domain.Article, error) {
	key := articleKey(articleOId)
	var article string
	if err := r.cache.Get(key, &article); err != nil {
		a, err := r.ArticleRepository.GetArticle(articleOId)
		if err == nil {
			str, e := r.serializer.Serialize([]domain.Article{a})
			if e != nil {
				return a, err
			}
			if err := r.cache.Set(key, str); err != nil {
				logrus.Warnln("failed to set cache:", err)
			}
		}
		return a, err
	}
	a, err := r.serializer.Deserialize(article)
	if err != nil {
		return domain.Article{}, err
	}
	return a[0], nil
}

func (r *CachedArticleRepository) GetForUser(username string, page, limit int) ([]domain.Article, error) {
	key := articlesKey(username, page, limit)
	var a string
	if err := r.cache.Get(key, &a); err != nil {
		articles, err := r.ArticleRepository.GetForUser(username, page, limit)
		if err == nil {
			str, e := r.serializer.Serialize(articles)
			if e != nil {
				return articles, err
			}
			if err := r.cache.Set(key, str); err != nil {
				logrus.Warnln("failed to set cache: ", err)
			}
		}
		return articles, err
	}

	return r.serializer.Deserialize(a)
}
