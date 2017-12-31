package inscription

import(
	"net/http"
	"database/sql"
	_"github.com/gorilla/sessions"
    "github.com/labstack/echo"
	_"github.com/labstack/echo/middleware"
	"crypto/md5"
	_"github.com/lib/pq"
	"encoding/hex"
	"strings"
	"fmt"
	"log"
)

func DisplayInscriptionPage(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "register.html", map[string]interface{}{})
}

func RegisterUser(c echo.Context, db *sql.DB) error {
	sqlStatement := "SELECT id FROM users WHERE username=$1"
	if c.FormValue("username") == "" {
		return c.Render(http.StatusNotAcceptable, "register.html", map[string]interface{}{
			"error": "username is required",
		})
	}
	if c.FormValue("firstname") == "" {
		return c.Render(http.StatusNotAcceptable, "register.html", map[string]interface{}{
			"error": "Firstname is required",
		})
	}
	if c.FormValue("lastname") == "" {
		return c.Render(http.StatusNotAcceptable, "register.html", map[string]interface{}{
			"error": "Lastname required",
		})
	}
	pos := strings.Index(c.FormValue("email"), "@")
	if strings.Contains(c.FormValue("email")[pos:], ".") == false {
		return c.Render(http.StatusNotAcceptable, "register.html", map[string]interface{}{
			"error": "Invalid mail address",
		})
	}
	rows, errs := db.Query(sqlStatement, c.FormValue("username"))
	if errs != nil {
		log.Println(errs)
		return c.Redirect(http.StatusMovedPermanently, "/error/1")
	}
	if (rows.Next()) {
		
		return c.Render(http.StatusNotAcceptable, "register.html", map[string]interface{}{
			"error": "Username already taken",
		})
	}
	sqlStatement = "SELECT id FROM users WHERE mail_adress=$1"
	rows, errs = db.Query(sqlStatement, c.FormValue("email"))
	if errs != nil {
		log.Println(errs)
		return c.Redirect(http.StatusMovedPermanently, "/error/1")            
	}
	if (rows.Next()) {
		return c.Render(http.StatusNotAcceptable, "register.html", map[string]interface{}{
			"error": "Email already taken",
		})
	}
	if len(c.FormValue("password")) < 6 {
		return c.Render(http.StatusNotAcceptable, "register.html", map[string]interface{}{
			"error": "Password too short",
		})
	}
	hasher := md5.New()
	hasher.Write([]byte(c.FormValue("password")))
	sqlStatement = "INSERT INTO users (username, password, mail_adress, birthdate, rights, first_name, last_name) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	res, errs := db.Query(sqlStatement, c.FormValue("username"), hex.EncodeToString(hasher.Sum(nil)), c.FormValue("email"), c.FormValue("birthdate"), "1", c.FormValue("firstname"), c.FormValue("lastname"))
	if errs != nil {
		log.Println(errs)
		return c.Redirect(http.StatusMovedPermanently, "/error/1")        
	} else {
		fmt.Println(res)
		return c.Render(http.StatusCreated, "login.html", map[string]interface{}{})
	}
	return c.Render(http.StatusOK, "ok", nil)
}