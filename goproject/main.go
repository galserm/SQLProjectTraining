package main

import (
    "net/http"
    "database/sql"
    "github.com/gorilla/sessions"
    "github.com/labstack/echo-contrib/session"
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    "fmt"
    "log"
    _"github.com/lib/pq"
    "html/template"
    "io"
    "crypto/md5"
    "encoding/hex"
    "strings"
)

type Post struct {
    ID int
    Content string
    UserID int
    LikesNumber int
    PicturePath string
    UserName string
}

type Comment struct {
    ID int
    Content string
    UserID int
    LikesNumber int
    PicturePath string
    UserName string
}

type Posts struct {
    Posts []Post
}
// template struct
type TemplateRenderer struct {
    templates *template.Template
}

//template rendering method
func (t *TemplateRenderer) Render(writer io.Writer, name string, data interface {}, ctx echo.Context) error {
    if viewContext, isMap := data.(map[string]interface{}); isMap {
        viewContext["reverse"] = ctx.Echo().Reverse
    }
    return t.templates.ExecuteTemplate(writer, name, data)
}

func main() {
    var err error
    db, err := sql.Open("postgres", "user=neotek password=kringstone dbname=db3 sslmode=disable")
    if err != nil {
        log.Fatal(err);
    }
    if err = db.Ping(); err != nil {
        panic(err)
    } else {
        fmt.Println("DB connected")
    }

    e := echo.New()
    e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
    renderer := &TemplateRenderer{
        templates: template.Must(template.ParseGlob("template/*.html")),
    }
    e.Renderer = renderer
    e.GET("/home", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        sess.Options = &sessions.Options {
            Path: "/",
            MaxAge: 86400 * 7,
            HttpOnly: true,
        }
        sess.Save(ctx.Request(), ctx.Response())
        if sess.Values["id"] != "" {
            return ctx.Redirect(http.StatusMovedPermanently, "/homes")
        }
        return ctx.Render(http.StatusOK, "login.html", map[string]interface{}{
            "name": "random",
        })
    }).Name = "data"

    e.POST("/inscription", func(c echo.Context) error {
        sqlStatement := "SELECT id FROM users WHERE username=$1"
        pos := strings.Index(c.FormValue("email"), "@")
        if strings.Contains(".", c.FormValue("email")[pos:]) == false {
            return c.Render(http.StatusOK, "register.html", map[string]interface{}{
                "error": "Invalid mail address",
            })
        }
        rows, _ := db.Query(sqlStatement, c.FormValue("username"))
        if (rows.Next()) {
            
            return c.Render(http.StatusOK, "register.html", map[string]interface{}{
                "error": "Username already taken",
            })
        }
        sqlStatement = "SELECT id FROM users WHERE mail_adress=$1"
        rows, err = db.Query(sqlStatement, c.FormValue("email"))
        if (rows.Next()) {
            return c.Render(http.StatusOK, "register.html", map[string]interface{}{
                "error": "Email already taken",
            })
        }
        hasher := md5.New()
        hasher.Write([]byte(c.FormValue("password")))
        sqlStatement = "INSERT INTO users (username, password, mail_adress, birthdate, rights, first_name, last_name) VALUES ($1, $2, $3, $4, $5, $6, $7)"
        res, errs := db.Query(sqlStatement, c.FormValue("username"), hex.EncodeToString(hasher.Sum(nil)), c.FormValue("email"), c.FormValue("birthdate"), "1", c.FormValue("firstname"), c.FormValue("lastname"))
        if errs != nil {
            fmt.Print("Fuck on")
            fmt.Println(errs)
        
        } else {
            fmt.Println(res)
            return c.Render(http.StatusOK, "login.html", map[string]interface{}{})
        }
        return c.Render(http.StatusOK, "ok", nil)
    }).Name = "Error"

    e.POST("/like_comment", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        sqlStatement := "SELECT id FROM comments_likes WHERE user_id=$1 AND comment_id=$2"
        res, _ := db.Query(sqlStatement, sess.Values["id"], ctx.FormValue("comment_id"))
        if res.Next() {
            sqlStatement = "SELECT likes_number FROM comments WHERE id=$1"
            res, _ = db.Query(sqlStatement, ctx.FormValue("comment_id"))
            var (
                likes_number int
            )
            for res.Next() {
                res.Scan(&likes_number)
            }
            likes_number = likes_number - 1
            sqlStatement = "UPDATE comments SET likes_number=$1 WHERE id=$2"
            db.Query(sqlStatement, likes_number, ctx.FormValue("comment_id"))
            sqlStatement = "DELETE FROM comments_likes WHERE user_id=$1 AND comment_id=$2"
            db.Query(sqlStatement, sess.Values["id"], ctx.FormValue("comment_id"))
            return ctx.Redirect(http.StatusMovedPermanently, "/post/" + ctx.FormValue("post_id"))            
        }        
        sqlStatement = "INSERT INTO comments_likes (user_id, comment_id) VALUES ($1, $2)"
        _, errs := db.Query(sqlStatement,  sess.Values["id"], ctx.FormValue("comment_id"))
        if errs != nil {
            ctx.Render(http.StatusOK, "home.html", map[string]interface{}{
                "error": "Oops, something went wrong, try again later",
            })
        } else {

            sqlStatement = "SELECT likes_number FROM comments WHERE id=$1"
            res, _ = db.Query(sqlStatement, ctx.FormValue("comment_id"))
            var (
                likes_number int
            )
            for res.Next() {
                res.Scan(&likes_number)
            }
            likes_number = likes_number + 1
            sqlStatement = "UPDATE comments SET likes_number=$1 WHERE id=$2"
            db.Query(sqlStatement, likes_number, ctx.FormValue("comment_id"))
        } 
        return ctx.Redirect(http.StatusMovedPermanently, "/post/" + ctx.FormValue("post_id"))
    })

    e.POST("/like_post", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        sqlStatement := "SELECT id FROM posts_likes WHERE user_id=$1 AND post_id=$2"
        res, _ := db.Query(sqlStatement, sess.Values["id"], ctx.FormValue("post_id"))
        if res.Next() {
            sqlStatement = "SELECT likes_number FROM posts WHERE id=$1"
            res, _ = db.Query(sqlStatement, ctx.FormValue("post_id"))
            var (
                likes_number int
            )
            for res.Next() {
                res.Scan(&likes_number)
            }
            likes_number = likes_number - 1
            sqlStatement = "UPDATE posts SET likes_number=$1 WHERE id=$2"
            db.Query(sqlStatement, likes_number, ctx.FormValue("post_id"))
            sqlStatement = "DELETE FROM posts_likes WHERE user_id=$1 AND post_id=$2"
            db.Query(sqlStatement, sess.Values["id"], ctx.FormValue("post_id"))
            if ctx.FormValue("origin") == "home" {
                return ctx.Redirect(http.StatusMovedPermanently, "/homes")
            }
            return ctx.Redirect(http.StatusMovedPermanently, "/post/" + ctx.FormValue("post_id"))            
        }        
        sqlStatement = "INSERT INTO posts_likes (user_id, post_id) VALUES ($1, $2)"
        _, errs := db.Query(sqlStatement,  sess.Values["id"], ctx.FormValue("post_id"))
        if errs != nil {
            ctx.Render(http.StatusOK, "home.html", map[string]interface{}{
                "error": "Oops, something went wrong, try again later",
            })
        } else {

            sqlStatement = "SELECT likes_number FROM posts WHERE id=$1"
            res, _ = db.Query(sqlStatement, ctx.FormValue("post_id"))
            var (
                likes_number int
            )
            for res.Next() {
                res.Scan(&likes_number)
            }
            likes_number = likes_number + 1
            sqlStatement = "UPDATE posts SET likes_number=$1 WHERE id=$2"
            db.Query(sqlStatement, likes_number, ctx.FormValue("post_id"))
        } 
        if ctx.FormValue("origin") == "home" {
            return ctx.Redirect(http.StatusMovedPermanently, "/homes")
        }
        return ctx.Redirect(http.StatusMovedPermanently, "/post/" + ctx.FormValue("post_id"))
    })

    e.GET("/inscription", func(ctx echo.Context) error {
        return ctx.Render(http.StatusOK, "register.html", map[string]interface{}{
            "name": "random",
        })
    }).Name = "Incription"

    e.GET("/logout", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        sess.Values["id"] = ""
        //sess.Save(ctx.Request(), ctx.Response())
        return ctx.Redirect(http.StatusMovedPermanently, "/login")
    })

    e.GET("/login", func(ctx echo.Context) error {
        return ctx.Render(http.StatusOK, "login.html", map[string]interface{}{})
    })

    e.POST("/connect", func(ctx echo.Context) error {
        var (
            id int
        )
        sqlStatement := "SELECT id FROM users WHERE username=$1 AND password=$2"
        hasher := md5.New()
        hasher.Write([]byte(ctx.FormValue("password")))
        res, _ := db.Query(sqlStatement, ctx.FormValue("username"), hex.EncodeToString(hasher.Sum(nil)))
        if (res.Next()) {
            res.Scan(&id)
            sess, _ := session.Get("session", ctx)
            sess.Values["id"] = id
            sess.Save(ctx.Request(), ctx.Response())
            ctx.Redirect(http.StatusMovedPermanently, "/homes")
        }
        return ctx.Render(http.StatusOK, "login.html", map[string]interface{}{
            "error": "Invalid username/password",
        })
    }).Name = "Result"
    
    e.POST("/add_comment", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        sqlStatement := "INSERT INTO comments (content, user_id, likes_number, picture_path, post_id) VALUES ($1, $2, $3, $4, $5)"
        
        _, err := db.Query(sqlStatement, ctx.FormValue("comment_content"), sess.Values["id"], "0", "", ctx.FormValue("post_id_comment"))
        if err != nil {
            return ctx.Redirect(http.StatusMovedPermanently, "/post/" + ctx.FormValue("post_id_comment"))
        } 
        return ctx.Redirect(http.StatusMovedPermanently, "/post/" + ctx.FormValue("post_id_comment"))    
    })

    e.GET("/homes", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        if sess.Values["id"] == "" {
            ctx.Redirect(http.StatusForbidden, "/connect")            
        }
        sqlStatement := "SELECT posts.id,posts.content,posts.user_id,posts.likes_number,posts.picture_path,users.username FROM posts JOIN users ON posts.user_id=users.id WHERE posts.user_id=$1"
        rows, _ := db.Query(sqlStatement, sess.Values["id"])
        type PostList []Post
        var result PostList
        for rows.Next() {
            post := Post{}
            rows.Scan(&post.ID, &post.Content, &post.UserID, &post.LikesNumber, &post.PicturePath, &post.UserName)
            result = append(result, post)
        }
        return ctx.Render(http.StatusOK, "home.html", map[string]interface{}{
            "posts": result,
        })
    }).Name = "Home"

    e.GET("/post/:id", func(ctx echo.Context) error {
        sqlStatement := "SELECT posts.id,posts.content,posts.user_id,posts.likes_number,posts.picture_path,users.username FROM posts JOIN users ON posts.user_id=users.id WHERE posts.id=$1"
        res, _ := db.Query(sqlStatement, ctx.Param("id"))
        post := Post{}        
        if res.Next() {
            res.Scan(&post.ID, &post.Content, &post.UserID, &post.LikesNumber, &post.PicturePath, &post.UserName)
        }
        sqlStatement = "SELECT comments.id,comments.content,comments.user_id,comments.likes_number,comments.picture_path,users.username FROM comments JOIN users ON comments.user_id=users.id WHERE comments.post_id=$1"
        res, _ = db.Query(sqlStatement, ctx.Param("id"))
        type CommentList []Comment
        var result CommentList
        for res.Next() {
            comment := Comment{}
            res.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.LikesNumber, &comment.PicturePath, &comment.UserName)
            result = append(result, comment)
        }
        return ctx.Render(http.StatusOK, "post.html", map[string]interface{}{
            "posts": post,
            "comments": result,
        })
    })

    e.GET("/update_post/:id", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        if sess.Values["id"] == "" {
            return ctx.Redirect(http.StatusMovedPermanently, "/")
        }
        sqlStatement := "SELECT id,content FROM posts WHERE id=$1"
        res, _ := db.Query(sqlStatement, ctx.Param("id"))
        post := Post{}
        if res.Next() {
            res.Scan(&post.ID, &post.Content)
        }
        return ctx.Render(http.StatusOK, "update_post.html", map[string]interface{}{
            "post": post,
        })
    }) 
    
    e.POST("/update_post", func(ctx echo.Context) error {
        sqlStatement := "UPDATE posts SET content=$1 WHERE id=$2"
        db.Query(sqlStatement, ctx.FormValue("content"), ctx.FormValue("post_id"))
        return ctx.Redirect(http.StatusMovedPermanently, "/post/" + ctx.FormValue("post_id"))
    })

    e.POST("/post", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        sqlStatement := "INSERT INTO posts (content, user_id, likes_number, picture_path) VALUES ($1, $2, $3, $4)"  
        _, errs := db.Query(sqlStatement, ctx.FormValue("content"), sess.Values["id"], "0", "")
        if errs != nil {
            return ctx.Render(http.StatusOK, "home.html", map[string]interface{}{
                "error": "Oops, something went wrong, please try again later",
            })
        } else {
            return ctx.Render(http.StatusOK, "home.html", map[string]interface{}{})
        }
    })

    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Logger.Fatal(e.Start(":8081"))
}