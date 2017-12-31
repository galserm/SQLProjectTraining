package main 

import (
	_"net/http"
	_"net/http/httptest"
    "database/sql"
    _"github.com/gorilla/sessions"
    _"github.com/labstack/echo-contrib/session"
    _"github.com/labstack/echo"
    _"github.com/labstack/echo/middleware"
	"fmt"
	"testing"
	"log"
	_"github.com/stretchr/testify/assert"
    _"github.com/lib/pq"
    _"html/template"
    _"crypto/md5"
    _"encoding/hex"
    _"strings"
    _"strconv"
    _"./render"
    _ "./connection"
    _"./inscription"
    _"./users"
    _"./posts"
    _"./comments"
)

func TestMain(t *testing.T) {
	var err error
	db, err := sql.Open("postgres", "user=neotek password=kringstone dbname=db3 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		panic(err) 
	} else {
		fmt.Println("DB connected")
	}
	//e := echo.New()
/*	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
    renderer := &render.TemplateRenderer{
        Templates: template.Must(template.ParseGlob("template/*.html")),
    }
	e.Renderer = renderer
	sess, _ := session.Get("session", e.AcquireContext())
	sess.Options = &sessions.Options {  
		HttpOnly: true,
	}*/
	//t.Error(connection.TestConnectUser(t));
}