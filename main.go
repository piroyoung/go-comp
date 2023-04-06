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

	r.Static("/assets", "./assets")

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

	r.GET("v1/stream", func(c *gin.Context) {
		p := c.Query("q")
		msg := make(chan string)
		go func() {
			defer close(msg)
			repo.Stream(c, p, 20, func(t string) error {
				msg <- t
				return nil
			})
		}()
		c.Stream(func(w io.Writer) bool {
			if m, ok := <-msg; ok {
				c.SSEvent("chunk", m)
				return true
			}
			return false
		})

	})
	r.Run()
}
