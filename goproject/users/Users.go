package users

import (
	"github.com/gorilla/sessions"
	"net/http"
	"database/sql"
	_"github.com/gorilla/sessions"
    "github.com/labstack/echo"
	_"github.com/labstack/echo/middleware"
	_"github.com/lib/pq"
	"log"
	"strings"
	"crypto/md5"
	"encoding/hex"
)

type User struct {
    ID int
    UserName string
    FirstName string
    LastName string
    MailAddress string
}


func DisplayUserPage(ctx echo.Context, db *sql.DB, sess *sessions.Session) error {
	sqlStatement := "SELECT id,username,first_name,last_name,mail_adress FROM users WHERE id=$1"
	rows, err := db.Query(sqlStatement, ctx.Param("id"))
	if err != nil {
		log.Println(err)
		return ctx.Redirect(http.StatusMovedPermanently, "/homes")
	}
	user := User{}
	if rows.Next() {
		rows.Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName, &user.MailAddress)
	}
	sqlStatement = "SELECT id FROM followers WHERE following_user_id=$1 and followed_user_id=$2"
	rows, err = db.Query(sqlStatement, sess.Values["id"], ctx.Param("id"))
	if err != nil {
		log.Println(err)
		return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
	}
	var action string
	action = "follow"
	if rows.Next() {
		action = "unfollow"           
	}
	return ctx.Render(http.StatusOK, "user.html", map[string]interface{}{
		"user": user,
		"action": action,
	})
}

func FollowUser(ctx echo.Context, db *sql.DB, sess *sessions.Session) error {
	var sqlStatement string
	if ctx.FormValue("method") == "follow" {
		sqlStatement = "INSERT INTO followers (following_user_id, followed_user_id) VALUES ($1, $2)"
	} else {
		sqlStatement = "DELETE FROM followers WHERE following_user_id=$1 AND followed_user_id=$2"
	}
	_, errs := db.Query(sqlStatement, sess.Values["id"], ctx.FormValue("id"))
	if errs != nil {
		return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
	}
	return ctx.Redirect(http.StatusMovedPermanently, "/users/" + ctx.FormValue("id"))
}

func GetFollowedUsers(ctx echo.Context, db *sql.DB, sess *sessions.Session) error {
	if sess.Values["id"] == "" {
		return ctx.Redirect(http.StatusForbidden, "/connect")
	}
	sqlStatement := "SELECT followers.followed_user_id,users.username FROM followers JOIN users ON followers.followed_user_id=users.id WHERE followers.following_user_id=$1"
	res, errs := db.Query(sqlStatement, sess.Values["id"])
	if errs != nil {
		return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
	}
	type UserList []User
	var result UserList
	for res.Next() {
		user := User{}
		res.Scan(&user.ID, &user.UserName)
		result = append(result, user)
	}
	return ctx.Render(http.StatusOK, "following.html", map[string]interface{}{
		"results": result,
	})
}

func GetUserUpdateForm(ctx echo.Context, db *sql.DB, sess *sessions.Session) error {
	if sess.Values["id"] == "" {
		ctx.Redirect(http.StatusForbidden, "/connect")
	}
	sqlStatement := "SELECT COUNT (followed_user_id) AS followingNumber FROM followers WHERE following_user_id=$1"
	res, errs := db.Query(sqlStatement, sess.Values["id"])
	if errs != nil {
		return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
	}
	var following_number int
	if res.Next() {
		res.Scan(&following_number)
	}
	return ctx.Render(http.StatusOK, "account.html", map[string]interface{}{
		"me": following_number,
	})
}

func UpdateUserInfos(ctx echo.Context, db *sql.DB, sess *sessions.Session) error {
	if (ctx.FormValue("username") != "") {
		sqlStatement := "UPDATE users SET username=$1 WHERE id=$2"
		_, errs := db.Query(sqlStatement, ctx.FormValue("username"), sess.Values["id"])
		if errs != nil {
			return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
		}
	}
	if (ctx.FormValue("password") != "") {
		if len(ctx.FormValue("password")) < 6 {
			return ctx.Render(http.StatusOK, "account.html", map[string]interface{}{
				"error": "Password too short",
			})
		}
		hasher := md5.New()
		hasher.Write([]byte(ctx.FormValue("password")))
		sqlStatement := "UPDATE users SET password=$1 WHERE id=$2"
		_, errs := db.Query(sqlStatement, hex.EncodeToString(hasher.Sum(nil)), sess.Values["id"])
		if errs != nil {
			return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
		}
	}
	if (ctx.FormValue("email") != "") {
		pos := strings.Index(ctx.FormValue("email"), "@")
		if strings.Contains(ctx.FormValue("email")[pos:], ".") == false {
			return ctx.Render(http.StatusOK, "account.html", map[string]interface{}{
				"error": "Invalid mail address",
			})
		}    
		sqlStatement := "UPDATE users SET mail_adress=$1 WHERE id=$2"
		_, errs := db.Query(sqlStatement, ctx.FormValue("email"), sess.Values["id"])
		if errs != nil {
			return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
		}
	}
	if (ctx.FormValue("firstname") != "") {
		sqlStatement := "UPDATE users SET first_name=$1 WHERE id=$2"
		_, errs := db.Query(sqlStatement, ctx.FormValue("firstname"), sess.Values["id"])
		if errs != nil {
			return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
		}
	}
	if (ctx.FormValue("lastname") != "") {
		sqlStatement := "UPDATE users SET last_name=$1 WHERE id=$2"
		_, errs := db.Query(sqlStatement, ctx.FormValue("lastname"), sess.Values["id"])
		if errs != nil {
			return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
		}
	}
	return ctx.Redirect(http.StatusMovedPermanently, "/homes")

}