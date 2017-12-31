package connection

import(
	"net/http"
	"database/sql"
	_"github.com/gorilla/sessions"
    "github.com/labstack/echo-contrib/session"
    "github.com/labstack/echo"
	_"github.com/labstack/echo/middleware"
	"crypto/md5"
	_"github.com/lib/pq"
	"encoding/hex"
)

func RenderLoginPage(ctx echo.Context, errorCode int) error {
	
	return ctx.Render(http.StatusOK, "login.html", map[string]interface{}{})
}

func ConnectUser(ctx echo.Context, db *sql.DB) error {
	var (
		id int
	)
	sqlStatement := "SELECT id FROM users WHERE username=$1 AND password=$2"
	hasher := md5.New()
	hasher.Write([]byte(ctx.FormValue("password")))
	res, errs := db.Query(sqlStatement, ctx.FormValue("username"), hex.EncodeToString(hasher.Sum(nil)))
	if errs != nil {
		return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
	}
	if (res.Next()) {
		res.Scan(&id)
		sess, _ := session.Get("session", ctx)
		sess.Values["id"] = id
		sess.Save(ctx.Request(), ctx.Response())
		return ctx.Redirect(http.StatusMovedPermanently, "/homes") 
	}
	return ctx.Render(http.StatusNotFound, "login.html", map[string]interface{}{
		"error": "Invalid username/password",
	})
}