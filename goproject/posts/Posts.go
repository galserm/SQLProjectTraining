package posts

import (
	"github.com/gorilla/sessions"
	"net/http"
	"database/sql"
	"github.com/labstack/echo-contrib/session"
	 "github.com/labstack/echo"
	_"github.com/labstack/echo/middleware"
	_"github.com/lib/pq"
	"log"
	"strconv"
)


type Post struct {
    ID int
    Content string
    UserID int
    LikesNumber int
    PicturePath string
    UserName string
    IsEditable bool
}

type Comment struct {
    ID int
    Content string
    UserID int
    LikesNumber int
    PicturePath string
    UserName string
    IsEditable bool
}

func LikePost(ctx echo.Context, db *sql.DB, sess *sessions.Session) error {

	sqlStatement := "SELECT id FROM posts_likes WHERE user_id=$1 AND post_id=$2"
	res, errs := db.Query(sqlStatement, sess.Values["id"], ctx.FormValue("post_id"))
	if errs != nil {
		return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
	}
	if res.Next() {
		sqlStatement = "SELECT likes_number FROM posts WHERE id=$1"
		res, errs = db.Query(sqlStatement, ctx.FormValue("post_id"))
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
		sqlStatement = "UPDATE posts SET likes_number=$1 WHERE id=$2"
		_, errs = db.Query(sqlStatement, likes_number, ctx.FormValue("post_id"))
		if errs != nil {
			return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
		}
		sqlStatement = "DELETE FROM posts_likes WHERE user_id=$1 AND post_id=$2"
		_, errs = db.Query(sqlStatement, sess.Values["id"], ctx.FormValue("post_id"))
		if errs != nil {
			return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
		}
		if ctx.FormValue("origin") == "home" {
			return ctx.Redirect(http.StatusMovedPermanently, "/homes")
		}
		return ctx.Redirect(http.StatusMovedPermanently, "/post/" + ctx.FormValue("post_id"))            
	}        
	sqlStatement = "INSERT INTO posts_likes (user_id, post_id) VALUES ($1, $2)"
	_, errs = db.Query(sqlStatement,  sess.Values["id"], ctx.FormValue("post_id"))
	if errs != nil {
		ctx.Render(http.StatusOK, "home.html", map[string]interface{}{
			"error": "Oops, something went wrong, try again later",
		})
	} else {

		sqlStatement = "SELECT likes_number FROM posts WHERE id=$1"
		res, errs = db.Query(sqlStatement, ctx.FormValue("post_id"))
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
		sqlStatement = "UPDATE posts SET likes_number=$1 WHERE id=$2"
		res, errs = db.Query(sqlStatement, likes_number, ctx.FormValue("post_id"))
		if errs != nil {
			return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
		}
	} 
	if ctx.FormValue("origin") == "home" {
		return ctx.Redirect(http.StatusMovedPermanently, "/homes")
	}
	return ctx.Redirect(http.StatusMovedPermanently, "/post/" + ctx.FormValue("post_id"))	
}

func GetPost(ctx echo.Context, db *sql.DB) error {
	sess, _ := session.Get("session", ctx)
	sqlStatement := "SELECT posts.id,posts.content,posts.user_id,posts.likes_number,posts.picture_path,users.username FROM posts JOIN users ON posts.user_id=users.id WHERE posts.id=$1"
	res, errs := db.Query(sqlStatement, ctx.Param("id"))
	if errs != nil {
		return ctx.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
            "error": "Bad Request",      
        })}
	post := Post{}        
	if res.Next() {
		res.Scan(&post.ID, &post.Content, &post.UserID, &post.LikesNumber, &post.PicturePath, &post.UserName)
		post_id := strconv.Itoa(post.UserID)
		session_id := strconv.Itoa(sess.Values["id"].(int))
		log.Println(post_id, sess.Values["id"])
		if  session_id == post_id {
			log.Println(post.IsEditable)                
			post.IsEditable = true
		}
		log.Println(post.IsEditable)
	} else {
		return ctx.Render(http.StatusNotFound, "error.html", map[string]interface{}{
            "error": "Post not found",      
        })
	}

	sqlStatement = "SELECT comments.id,comments.content,comments.user_id,comments.likes_number,comments.picture_path,users.username FROM comments JOIN users ON comments.user_id=users.id WHERE comments.post_id=$1 ORDER BY comments.comment_date ASC"
	res, errs = db.Query(sqlStatement, ctx.Param("id"))
	if errs != nil {
		return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
	}
	type CommentList []Comment
	var result CommentList
	for res.Next() {
		comment := Comment{}
		res.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.LikesNumber, &comment.PicturePath, &comment.UserName)
		comment_id := strconv.Itoa(comment.UserID)
		session_id := strconv.Itoa(sess.Values["id"].(int))
		if session_id == comment_id {
			comment.IsEditable = true
		}
		result = append(result, comment)
		
	}
	return ctx.Render(http.StatusOK, "post.html", map[string]interface{}{
		"posts": post,
		"comments": result,
	})
}

func GetUpdatePostForm(ctx echo.Context, db *sql.DB, sess *sessions.Session) error {
	if sess.Values["id"] == "" {
		return ctx.Redirect(http.StatusMovedPermanently, "/")
	}
	sqlStatement := "SELECT id,content,user_id FROM posts WHERE id=$1"
	res, errs := db.Query(sqlStatement, ctx.Param("id"))
	if errs != nil {
		return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
	}
	post := Post{}
	if res.Next() {
		res.Scan(&post.ID, &post.Content, &post.UserID)
	}
	if post.UserID != sess.Values["id"] {
		return ctx.Redirect(http.StatusMovedPermanently, "/")
	}
	return ctx.Render(http.StatusOK, "update_post.html", map[string]interface{}{
		"post": post,
	})
}

func UpdatePost(ctx echo.Context, db *sql.DB) error {
	sqlStatement := "UPDATE posts SET content=$1 WHERE id=$2"
	if ctx.FormValue("content") == "" {
		return ctx.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
            "error": "BadRequest",      
        })
	}
	_, errs := db.Query(sqlStatement, ctx.FormValue("content"), ctx.FormValue("post_id"))
	if errs != nil {
		return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
	}
	return ctx.Redirect(http.StatusMovedPermanently, "/post/" + ctx.FormValue("post_id"))
}

func CreatePost(ctx echo.Context, db *sql.DB, sess *sessions.Session) error {
	sqlStatement := "INSERT INTO posts (content, user_id, likes_number, picture_path) VALUES ($1, $2, $3, $4)"  
	if ctx.FormValue("content") != "" {
		_, errs := db.Query(sqlStatement, ctx.FormValue("content"), sess.Values["id"], "0", "")
		if errs != nil {
			return ctx.Render(http.StatusOK, "home.html", map[string]interface{}{
				"error": "Oops, something went wrong, please try again later",
			})
		}
	} else {
		return ctx.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
            "error": "BadRequest",      
        })
	}	
	return ctx.Redirect(http.StatusMovedPermanently, "/homes")
}

func DeletePost(ctx echo.Context, db *sql.DB) error {
	sqlStatement := "DELETE FROM posts WHERE id=$1"
	_, errs := db.Query(sqlStatement, ctx.FormValue("post_id"))
	if errs != nil {
		return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
	}
	return ctx.Redirect(http.StatusMovedPermanently, "/homes")
}