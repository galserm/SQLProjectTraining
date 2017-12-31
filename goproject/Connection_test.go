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
	connection "./connection"
	_"github.com/labstack/echo-contrib/session"
	_"github.com/gorilla/sessions"
	"./render"
	"html/template"
)

func TestConnectUserWithGoodLogs(t *testing.T) {
	e := echo.New()
	db, _:= sql.Open("postgres", "user=neotek password=kringstone dbname=db3 sslmode=disable")	
	form := make(url.Values)
	form.Set("password", "kringstone")
	form.Set("username", "n30t3k")
	req := httptest.NewRequest("POST", "/connect", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
    renderer := &render.TemplateRenderer{
        Templates: template.Must(template.ParseGlob("template/*.html")),
    }
	e.Renderer = renderer
	
	ctx := e.NewContext(req, rec)
	if assert.NoError(t, connection.ConnectUser(ctx, db)){
		assert.Equal(t, http.StatusMovedPermanently, rec.Code)
	}
}

func TestConnectUserWithBadLogs(t *testing.T) {
	e := echo.New()
	db, _:= sql.Open("postgres", "user=neotek password=kringstone dbname=db3 sslmode=disable")	
	form := make(url.Values)
	form.Set("password", "kringstone")
	form.Set("username", "n30tek")
	req := httptest.NewRequest("POST", "/connect", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
    renderer := &render.TemplateRenderer{
        Templates: template.Must(template.ParseGlob("template/*.html")),
    }
	e.Renderer = renderer
	
	ctx := e.NewContext(req, rec)
	
	if assert.NoError(t, connection.ConnectUser(ctx, db)){
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
	
}