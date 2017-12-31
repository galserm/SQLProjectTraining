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
    _"crypto/md5"
    _"encoding/hex"
    _"strings"
    "strconv"
    "./render"
    "./connection"
    "./inscription"
    "./users"
    "./posts"
    "./comments"
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

type User struct {
    ID int
    UserName string
    FirstName string
    LastName string
    MailAddress string
}

type Posts struct {
    Posts []Post
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
    renderer := &render.TemplateRenderer{
        Templates: template.Must(template.ParseGlob("template/*.html")),
    }
    e.Renderer = renderer
     
    e.GET("/", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        sess.Options = &sessions.Options {  
            HttpOnly: true,
        }
        sess.Save(ctx.Request(), ctx.Response())
        if _, ok := sess.Values["id"]; ok {
            return ctx.Redirect(http.StatusMovedPermanently, "/login")
        }
        return ctx.Render(http.StatusOK, "login.html", map[string]interface{}{
            "name": "random",
        })
    }).Name = "data"

    e.GET("/error/:code", func(ctx echo.Context) error {
        errmess := ""
        errcode := ctx.Param("code")
        if errcode == "1" {
            errmess = "Unknown error, try again later"
        }
        return ctx.Render(http.StatusBadRequest, "error.html", map[string]interface{}{
            "error": errmess,      
        })
    })

    e.POST("/inscription", func(c echo.Context) error {
        return inscription.RegisterUser(c, db)
    }).Name = "Error"

    e.POST("/like_comment", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        return comments.LikeComment(ctx, db, sess)
    })

    e.POST("/like_post", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        return posts.LikePost(ctx, db, sess)
    })

    e.GET("/inscription", func(ctx echo.Context) error {
        return inscription.DisplayInscriptionPage(ctx)
    }).Name = "Incription"

    e.GET("/logout", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        delete(sess.Values, "id")
        sess.Save(ctx.Request(), ctx.Response())
        return ctx.Redirect(http.StatusMovedPermanently, "/login")
    })

    e.GET("/login", func(ctx echo.Context) error {
        return connection.RenderLoginPage(ctx, http.StatusOK)
    })

    e.POST("/connect", func(ctx echo.Context) error {
        return connection.ConnectUser(ctx, db);
     }).Name = "Result"
    
    e.GET("/users/:id", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        return users.DisplayUserPage(ctx, db, sess)
    })

    e.POST("/follow", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        return users.FollowUser(ctx, db, sess)
    })

    e.POST("/add_comment", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        return comments.CreateComment(ctx, db, sess)    
    })

    e.GET("/homes", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        if _, ok := sess.Values["id"]; !ok {
            return ctx.Redirect(http.StatusMovedPermanently, "/login")
        }
        sqlStatement := "SELECT posts.id,posts.content,posts.user_id,posts.likes_number,posts.picture_path,users.username FROM posts JOIN users ON posts.user_id=users.id WHERE posts.user_id=$1 ORDER BY posts.post_date ASC"
        rows, errs := db.Query(sqlStatement, sess.Values["id"])
        if errs != nil {
            return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
        }
        type PostList []Post
        var result PostList
        for rows.Next() {
            post := Post{}
            rows.Scan(&post.ID, &post.Content, &post.UserID, &post.LikesNumber, &post.PicturePath, &post.UserName)
            result = append(result, post)
        }
        sqlStatement = "SELECT followed_user_id FROM followers WHERE following_user_id=$1"
        rows, errs = db.Query(sqlStatement, sess.Values["id"])
        if errs != nil {
            return ctx.Redirect(http.StatusMovedPermanently, "/error/1")    
        }
        var following_user_id []int
        var user_id int
        for rows.Next() {
            rows.Scan(&user_id)
            following_user_id = append(following_user_id, user_id)   
        }
        for i := 0; i < len(following_user_id); i += 1 {
            sqlStatement = "SELECT posts.id,posts.content,posts.user_id,posts.likes_number,posts.picture_path,users.username FROM posts JOIN users ON posts.user_id=users.id WHERE posts.user_id=$1"
            rows, errs = db.Query(sqlStatement, following_user_id[i])
            if errs != nil {
                return ctx.Redirect(http.StatusMovedPermanently, "/error/1")
            }
            for rows.Next() {
                newpost := Post{}
                rows.Scan(&newpost.ID, &newpost.Content, &newpost.UserID, &newpost.LikesNumber, &newpost.PicturePath, &newpost.UserName)
                result = append(result, newpost)
            }
        }
        return ctx.Render(http.StatusOK, "home.html", map[string]interface{}{
            "posts": result,
        })
    }).Name = "Home"

    e.GET("/post/:id", func(ctx echo.Context) error {
        return posts.GetPost(ctx, db)
    })

    e.POST("/delete_comment", func(ctx echo.Context) error {
        return comments.DeleteComment(ctx, db)
    })

    e.POST("/delete_post", func(ctx echo.Context) error {
        return posts.DeletePost(ctx, db)
    })

    e.GET("/search", func(ctx echo.Context) error {
        return ctx.Render(http.StatusOK, "search.html", map[string]interface{}{})
    }) 

    e.POST("/search", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        sqlStatement := "SELECT DISTINCT id,username FROM users WHERE username LIKE '%" + ctx.FormValue("pattern") + "%'"
        rows, err := db.Query(sqlStatement)
        if err != nil {
            log.Println(err)
            return ctx.Redirect(http.StatusMovedPermanently, "/search")
        }
        type UserList []User
        var result UserList
        var i int = 0
        log.Println(rows)
        for rows.Next() {
            user := User{}
            rows.Scan(&user.ID, &user.UserName)
            if sess.Values["id"] != strconv.Itoa(user.ID) {
                result = append(result, user)
            }
            i = i + 1
        }
        if i == 0 {
            return ctx.Render(http.StatusOK, "search.html", map[string]interface{}{
                "error": "No results",
            })
        }
        return ctx.Render(http.StatusOK, "search.html", map[string]interface{}{
            "results": result,
        })
    }) 

    e.GET("/following", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        return users.GetFollowedUsers(ctx, db, sess)
    })

    e.GET("/user_update", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        return users.GetUserUpdateForm(ctx, db, sess)
    })

    e.POST("/user_update", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        return users.UpdateUserInfos(ctx, db, sess)
    })

    e.GET("/update_comment/:id", func (ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        return comments.GetCommentUpdateForm(ctx, db, sess)
    }) 

    e.POST("/update_comment", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        return comments.UpdateComment(ctx, db, sess)
    })

    e.GET("/update_post/:id", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        return posts.GetUpdatePostForm(ctx, db, sess)
    }) 
    
    e.POST("/update_post", func(ctx echo.Context) error {
        return posts.UpdatePost(ctx, db)
    })

    e.POST("/post", func(ctx echo.Context) error {
        sess, _ := session.Get("session", ctx)
        return posts.CreatePost(ctx, db, sess)
    })

    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Logger.Fatal(e.Start(":8081"))
} 