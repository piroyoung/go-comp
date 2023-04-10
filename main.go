package main

import (
	"encoding/json"
	"fmt"
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

	r.POST("v1/stream", func(c *gin.Context) {
		var prompt Prompt
		if err := c.BindJSON(&prompt); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		msg := make(chan string)
		go func() {
			defer close(msg)
			repo.Stream(c, prompt.Value, 200, func(t string) error {
				msg <- t
				return nil
			})
		}()
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Stream(func(w io.Writer) bool {
			if m, ok := <-msg; ok {
				fmt.Print(m)
				jsonBytes, _ := json.Marshal(Token{Value: m})
				c.SSEvent("chunk", string(jsonBytes))
				return true
			}
			return false
		})

	})
	r.Run()
}
