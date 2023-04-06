package main

import (
	"github.com/gin-gonic/gin"
	"github.com/piroyoung/go-aoai"
	"github.com/piroyoung/go-comp/repository"
	"io"
	"os"
)

func main() {
	resourceName := "example-aoai-02"
	deploymentName := "gpt-35-turbo-0301"
	apiVersion := "2023-03-15-preview"
	accessToken := os.Getenv("AZURE_OPENAI_API_KEY")

	client := aoai.New(resourceName, deploymentName, apiVersion, accessToken)
	repo := repository.NewCompletionRepository(client)

	r := gin.Default()

	r.GET("v1/completion", func(c *gin.Context) {
		t, err := repo.Complete(c, c.Query("q"), 20)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"text": t,
		})
	})

	r.GET("v1/stream/:q", func(c *gin.Context) {
		q := c.Param("q")
		c.Stream(func(w io.Writer) bool {
			t, err := repo.Complete(c, q, 20)
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return false
			}
			c.SSEvent("message", t)
			return true
		})
	})

	r.Run()
}
