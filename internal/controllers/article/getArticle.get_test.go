package article

import (
	"errors"
	"net/http"
	"testing"
	"user-management/internal/domain"
	cntlrMocks "user-management/internal/mocks/controllers"
	mdlMocks "user-management/internal/mocks/models"
	"user-management/internal/models"
	"user-management/internal/views"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestArticleController_Get(t *testing.T) {
	type mockedRes struct {
		article       domain.Article
		getErr        error
		jsonCode      int
		jsonBody      interface{}
		noContentCode int
	}
	ctxMock := cntlrMocks.NewContext(t)
	articleModelMock := mdlMocks.NewArticleModel(t)
	setup := func(res mockedRes) func() {
		articleId := "artice-id"
		pathParamCall := ctxMock.EXPECT().PathParam("article_id").Return(articleId).Once()
		getCall := articleModelMock.EXPECT().Get(articleId).NotBefore(pathParamCall).Return(res.article, res.getErr).Once()
		calls := []*mock.Call{
			pathParamCall, getCall,
			ctxMock.EXPECT().JSON(res.jsonCode, res.jsonBody).NotBefore(pathParamCall, getCall).Return(nil).Maybe(),
			ctxMock.EXPECT().NoContent(res.noContentCode).NotBefore(pathParamCall, getCall).Return(nil).Maybe(),
		}
		return func() {
			for _, call := range calls {
				call.Unset()
			}
		}
	}
	someError := errors.New("some error")
	artcl := domain.Article{Theme: "theme", Text: "text", Tags: []string{}}
	tests := []struct {
		name      string
		mockedRes mockedRes
		wantErr   error
	}{
		{
			name: "article is not found",
			mockedRes: mockedRes{
				getErr:   models.NotFoundErr,
				jsonCode: http.StatusNotFound,
				jsonBody: mock.Anything,
			},
			wantErr: models.NotFoundErr,
		},
		{
			name: "internal server error",
			mockedRes: mockedRes{
				getErr:        someError,
				noContentCode: http.StatusInternalServerError,
			},
			wantErr: someError,
		},
		{
			name: "success",
			mockedRes: mockedRes{
				article:  artcl,
				jsonCode: http.StatusOK,
				jsonBody: views.NewArticleView(artcl),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanSetup := setup(tt.mockedRes)
			defer cleanSetup()
			err := NewArticleController(articleModelMock).Get(ctxMock)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}
