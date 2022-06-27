package posts

import (
	"context"
	"fmt"
	"strconv"

	// "github.com/JesusJMM/blog-plat-go/postgres"
	"github.com/gin-gonic/gin"
	"github.com/vingarcia/ksql"
)

type PostsHandler struct {
	db  *ksql.DB
	ctx context.Context
}

func New(db *ksql.DB, ctx context.Context) PostsHandler {
	return PostsHandler{
		db:  db,
		ctx: ctx,
	}
}

type PartialPostWithAuthor struct {
	Article PartialArticle `tablename:"a"`
	Author  Author         `tablename:"u"`
}

func (h PostsHandler) AllArticles() gin.HandlerFunc {
	return func(c *gin.Context) {
		var posts []PartialPostWithAuthor
		err := h.db.Query(h.ctx, &posts,
			"FROM articles as a LEFT JOIN users as u on a.user_id = u.user_id",
		)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"posts": posts})
	}
}

func (h PostsHandler) ArticlesPaginated() gin.HandlerFunc {
	return func(c *gin.Context) {
		queryPage := c.DefaultQuery("page", "1")
		page, err := strconv.Atoi(queryPage)
		if err != nil {
			c.JSON(500, gin.H{"error": "'page' query param must be a number"})
			return
		}
		var posts []PartialPostWithAuthor
    q := fmt.Sprintf(`
    FROM articles as a
    LEFT JOIN users as u
    ON a.user_id = u.user_id
    ORDER BY a.article_id
    LIMIT %d
    OFFSET %d
    `, 10, (page -1) * 10)
		err = h.db.Query(
      h.ctx, 
      &posts,
      q,
		)
    if err != nil {
      c.JSON(500, gin.H{"error": err.Error()})
      return
    }
    c.JSON(200, gin.H{"posts": posts})
	}
}
