package comments

import (
	"github.com/gorilla/sessions"
	"net/http"
	"database/sql"
	_"github.com/gorilla/sessions"
    "github.com/labstack/echo"
	_"github.com/labstack/echo/middleware"
	_"github.com/lib/pq"
	_"log"
	_"strconv"
)

type Comment struct {
    ID int
    Content string
    UserID int
    LikesNumber int
    PicturePath string
    UserName string
    IsEditable bool
}

func LikeComment(ctx echo.Context, db *sql.DB, sess *sessions.Session) error {
	sqlStatement := "SELECT id FROM comments_likes WHERE user_id=$1 AND comment_id=$2"
	res, errs := db.Query(sqlStatement, sess.Values["id"], ctx.FormValue("comment_id"))
	if errs != nil {
		return ctx.Redirect(http.StatusMovedPermanently, "/error/1")            
	}
	if res.Next() {
		sqlStatement = "SELECT likes_number FROM comments WHERE id=$1"
		res, errs = db.Query(sqlStatement, ctx.FormValue("comment_id"))
		if errs != nil {
			return ctx.Redirect(http.StatusMovedPermanently, "/error/1")            
		}
		var (
			likes_number int
		)
		for res.Next() {
			res.Scan(&likes_number)
		}
		likes_number = likes_number - 1
		sqlStatement = "UPDATE comments SET likes_number=$1 WHERE id=$2"
		_, errs = db.Query(sqlStatement, likes_number, ctx.FormValue("comment_id"))
		if errs != nil {
			return ctx.Redirect(http.StatusMovedPermanently, "/error/1")        
		}
		sqlStatement = "DELETE FROM comments_likes WHERE user_id=$1 AND comment_id=$2"
		_, errs = db.Query(sqlStatement, sess.Values["id"], ctx.FormValue("comment_id"))
		if errs != nil {
			return ctx.Redirect(http.StatusMovedPermanently, "/error/1")        
		}
		return ctx.Redirect(http.StatusMovedPermanently, "/post/" + ctx.FormValue("post_id"))            
	}        
	sqlStatement = "INSERT INTO comments_likes (user_id, comment_id) VALUES ($1, $2)"
	_, errs = db.Query(sqlStatement,  sess.Values["id"], ctx.FormValue("comment_id"))
	if errs != nil {
		ctx.Render(http.StatusOK, "home.html", map[string]interface{}{
			"error": "Oops, something went wrong, try again later",
		})
	} else {

		sqlStatement = "SELECT likes_number FROM comments WHERE id=$1"
		res, errs = db.Query(sqlStatement, ctx.FormValue("comment_id"))
		if errs != nil {
			return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
		}
		var (
			likes_number int
		)
		for res.Next() {
			res.Scan(&likes_number)
		}
		likes_number = likes_number + 1
		sqlStatement = "UPDATE comments SET likes_number=$1 WHERE id=$2"
		_, errs = db.Query(sqlStatement, likes_number, ctx.FormValue("comment_id"))
		if errs != nil {
			return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
		}
	} 
	return ctx.Redirect(http.StatusMovedPermanently, "/post/" + ctx.FormValue("post_id"))
}

func CreateComment(ctx echo.Context, db *sql.DB, sess *sessions.Session) error {
	sqlStatement := "INSERT INTO comments (content, user_id, likes_number, picture_path, post_id) VALUES ($1, $2, $3, $4, $5)"
	
	if ctx.FormValue("comment_content") != "" {
		_, err := db.Query(sqlStatement, ctx.FormValue("comment_content"), sess.Values["id"], "0", "", ctx.FormValue("post_id_comment"))
		if err != nil {
			return ctx.Redirect(http.StatusMovedPermanently, "/post/" + ctx.FormValue("post_id_comment"))
		} 
	}
	return ctx.Redirect(http.StatusMovedPermanently, "/post/" + ctx.FormValue("post_id_comment"))
}

func GetCommentUpdateForm(ctx echo.Context, db *sql.DB, sess *sessions.Session) error {
	if sess.Values["id"] == "" {
		return ctx.Redirect(http.StatusMovedPermanently, "/login")
	}
	sqlStatement := "SELECT id,content,user_id FROM comments WHERE id=$1"
	res, errs := db.Query(sqlStatement, ctx.Param("id"))
	if errs != nil {
		return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
	}
	comment := Comment{}
	if res.Next() {
		res.Scan(&comment.ID, &comment.Content, &comment.UserID)
	}
	if comment.UserID != sess.Values["id"] {
		return ctx.Redirect(http.StatusMovedPermanently, "/")
	}
	return ctx.Render(http.StatusOK, "update_comment.html", map[string]interface{}{
		"comment": comment,
	})
}

func UpdateComment(ctx echo.Context, db *sql.DB, sess *sessions.Session) error {
	if sess.Values["id"] == "" {
		return ctx.Redirect(http.StatusMovedPermanently, "/")
	}
	sqlStatement := "UPDATE comments SET content=$1 WHERE id=$2"
	_, errs := db.Query(sqlStatement, ctx.FormValue("content"), ctx.FormValue("comment_id"))
	if errs != nil {
		return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
	}
	return ctx.Redirect(http.StatusMovedPermanently, "/update_comment/" + ctx.FormValue("comment_id"))
}

func DeleteComment(ctx echo.Context, db *sql.DB) error {
	sqlStatement := "DELETE FROM comments WHERE id=$1"
	_, errs := db.Query(sqlStatement, ctx.FormValue("comment_id"))
	if errs != nil {
		return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
	}
	return ctx.Redirect(http.StatusMovedPermanently, "/homes")
}
