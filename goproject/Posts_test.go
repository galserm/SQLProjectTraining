package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"net/url"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"database/sql"
	_"github.com/lib/pq"
	posts "./posts"
	_"github.com/labstack/echo-contrib/session"
	_"github.com/gorilla/sessions"
	"./render"
	"html/template"
)

/func TestGetPostWithValidId(t *testing.T) {
	e := echo.New()
	db, _:= sql.Open("postgres", "user=neotek password=kringstone dbname=db3 sslmode=disable")	
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
    renderer := &render.TemplateRenderer{
        Templates: template.Must(template.ParseGlob("template/*.html")),
    }
	e.Renderer = renderer
	
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/posts/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("40")
	if assert.NoError(t, posts.GetPost(ctx, db)){
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestGetPostWithNotValidId(t *testing.T) {
	e := echo.New()
	db, _:= sql.Open("postgres", "user=neotek password=kringstone dbname=db3 sslmode=disable")	
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
    renderer := &render.TemplateRenderer{
        Templates: template.Must(template.ParseGlob("template/*.html")),
    }
	e.Renderer = renderer
	
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/posts/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("40000")
	if assert.NoError(t, posts.GetPost(ctx, db)){
		assert.Equal(t, http.StatusOK, rec.Code)
	}	
}

func TestUpdatePostWithValidInput(t *testing.T) {
	e := echo.New()
	db, _:= sql.Open("postgres", "user=neotek password=kringstone dbname=db3 sslmode=disable")	
	form := make(url.Values)
	form.Set("content", "test_updatedcontent")
	form.Set("post_id", "40")
	req := httptest.NewRequest("POST", "/inscription", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
    renderer := &render.TemplateRenderer{
        Templates: template.Must(template.ParseGlob("template/*.html")),
    }
	e.Renderer = renderer
	
	ctx := e.NewContext(req, rec)
	if assert.NoError(t, posts.UpdatePost(ctx, db)){
		assert.Equal(t, http.StatusMovedPermanently, rec.Code)
	}
}

func TestUpdatePostWithNotValidInput(t *testing.T) {
	e := echo.New()
	db, _:= sql.Open("postgres", "user=neotek password=kringstone dbname=db3 sslmode=disable")	
	form := make(url.Values)
	form.Set("content", "")
	form.Set("post_id", "40")
	req := httptest.NewRequest("POST", "/update_post", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
    renderer := &render.TemplateRenderer{
        Templates: template.Must(template.ParseGlob("template/*.html")),
    }
	e.Renderer = renderer
	
	ctx := e.NewContext(req, rec)
	if assert.NoError(t, posts.UpdatePost(ctx, db)){
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestCreatePostWithValidInput(t *testing.T) {
	e := echo.New()
	db, _:= sql.Open("postgres", "user=neotek password=kringstone dbname=db3 sslmode=disable")	
	form := make(url.Values)
	form.Set("content", "test_content")
	req := httptest.NewRequest("POST", "/post", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
    renderer := &render.TemplateRenderer{
        Templates: template.Must(template.ParseGlob("template/*.html")),
    }
	e.Renderer = renderer
	
	ctx := e.NewContext(req, rec)
 	if assert.NoError(t, posts.UpdatePost(ctx, db)){
		assert.Equal(t, http.StatusMovedPermanently, rec.Code)
	}
}

func TestCreatePostWithNotValidInput(t *testing.T) {
	e := echo.New()
	db, _:= sql.Open("postgres", "user=neotek password=kringstone dbname=db3 sslmode=disable")	
	form := make(url.Values)
	form.Set("content", "")
	req := httptest.NewRequest("POST", "/post", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
    renderer := &render.TemplateRenderer{
        Templates: template.Must(template.ParseGlob("template/*.html")),
    }
	e.Renderer = renderer
	
	ctx := e.NewContext(req, rec)
 	if assert.NoError(t, posts.UpdatePost(ctx, db)){
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}