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
	inscription "./inscription"
	_"github.com/labstack/echo-contrib/session"
	_"github.com/gorilla/sessions"
	"./render"
	"html/template"
)

/*func TestRegisterUserWithValidInput(t *testing.T) {
	e := echo.New()
	db, _:= sql.Open("postgres", "user=neotek password=kringstone dbname=db3 sslmode=disable")	
	form := make(url.Values)
	form.Set("password", "test_password")
	form.Set("username", "test_username")
	form.Set("email", "testemail@email.com")
	form.Set("birthdate", "2017-12-09 00:00:00+08")
	form.Set("firstname", "test_firstname")
	form.Set("lastname", "test_lastname")
	req := httptest.NewRequest("POST", "/inscription", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
    renderer := &render.TemplateRenderer{
        Templates: template.Must(template.ParseGlob("template/*.html")),
    }
	e.Renderer = renderer
	
	ctx := e.NewContext(req, rec)
	if assert.NoError(t, inscription.RegisterUser(ctx, db)){
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}*/

func TestRegisterUserWithAlreadyExistingUsername(t *testing.T) {
	e := echo.New()
	db, _:= sql.Open("postgres", "user=neotek password=kringstone dbname=db3 sslmode=disable")	
	form := make(url.Values)
	form.Set("password", "test_password")
	form.Set("username", "test_username")
	form.Set("email", "testemail@email.com")
	form.Set("birthdate", "2017-12-09 00:00:00+08")
	form.Set("firstname", "test_firstname")
	form.Set("lastname", "test_lastname")
	req := httptest.NewRequest("POST", "/inscription", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
    renderer := &render.TemplateRenderer{
        Templates: template.Must(template.ParseGlob("template/*.html")),
    }
	e.Renderer = renderer
	
	ctx := e.NewContext(req, rec)
	if assert.NoError(t, inscription.RegisterUser(ctx, db)){
		assert.Equal(t, http.StatusNotAcceptable, rec.Code)
	}
}