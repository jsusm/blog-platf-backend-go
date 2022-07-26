package api

import (
	"context"

	"github.com/JesusJMM/blog-plat-go/api/articles"
	"github.com/JesusJMM/blog-plat-go/api/authors"
	"github.com/JesusJMM/blog-plat-go/api/auth"
	"github.com/JesusJMM/blog-plat-go/postgres/repos/users"
	articlesRepo "github.com/JesusJMM/blog-plat-go/postgres/repos/articles"
	"github.com/gin-gonic/gin"
  "github.com/gin-contrib/cors"
	"github.com/vingarcia/ksql"
)

// Return a instance of gin.Engine with all routes
// registered
func New(db *ksql.DB) gin.Engine {
	r := gin.Default()
  corsConfig := cors.DefaultConfig()
  corsConfig.AllowAllOrigins = true 
  r.Use(cors.New(corsConfig))
	api := r.Group("/api")

	articleH := articles.New(db, context.Background(), articlesRepo.NewArticleRepo(db, context.Background()))
	authH := auth.New(db, context.Background(), users.New(db, context.Background()))
  authorsC := authors.New(db, context.Background())

	api.GET("/articles/all", articleH.All())
	api.GET("/articles/paginated", articleH.Paginated())
	api.GET("/articles/author/:author", articleH.ByAuthorPaginated())
  api.GET("/article/:author/:slug", articleH.OneArticle())

  api.POST("/article", auth.AuthRequired, articleH.Create())
  api.PUT("/article", auth.AuthRequired, articleH.Update())
  api.DELETE("/article/:id", auth.AuthRequired, articleH.Delete())

	api.POST("/auth/signup", authH.Signup())
	api.POST("/auth/login", authH.Login())

  api.GET("/authors/all", authorsC.GetAll())
  api.GET("/author/:authorName", authorsC.One())
	return *r
}
