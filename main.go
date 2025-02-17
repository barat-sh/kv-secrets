package main

import (
	"log"
	"os"
  "fmt"
  "net/http"
	"github.com/joho/godotenv"
  "github.com/gin-gonic/gin"

	"github.com/go-resty/resty/v2"
)

var client = resty.New()

func writeToKv(c *gin.Context){
  key := c.Param("key")
  value := c.PostForm("value")
  connectionURL, exists := os.LookupEnv("KV_CONNECTION_URL")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing KV_CONNECTION_URL"})
		return
	}

	accountID, exists := os.LookupEnv("KV_ACCOUNT_ID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing KV_ACCOUNT_ID"})
		return
	}

	namespaceID, exists := os.LookupEnv("KV_NAMESPACE_ID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing KV_NAMESPACE_ID"})
		return
	}

	apiToken, exists := os.LookupEnv("KV_API_TOKEN")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing KV_API_TOKEN"})
		return
	} 
  fmt.Println(key)
  fmt.Println(value)
  url := fmt.Sprintf(connectionURL, accountID, namespaceID, key)

  resp, err := client.R().SetAuthToken(apiToken).SetHeader("Content-Type", "text/plain").SetBody(value).Put(url)

  if err != nil{
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }

  c.Data(resp.StatusCode(), "application/json", []byte(resp.String()))
}

func readFromKv(c *gin.Context) {
  key := c.Param("key")
  connectionURL, exists := os.LookupEnv("KV_CONNECTION_URL")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing KV_CONNECTION_URL"})
		return
	}

	accountID, exists := os.LookupEnv("KV_ACCOUNT_ID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing KV_ACCOUNT_ID"})
		return
	}

	namespaceID, exists := os.LookupEnv("KV_NAMESPACE_ID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing KV_NAMESPACE_ID"})
		return
	}

	apiToken, exists := os.LookupEnv("KV_API_TOKEN")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing KV_API_TOKEN"})
		return
	}
  url := fmt.Sprintf(connectionURL, accountID, namespaceID, key)

  resp, err := client.R().SetAuthToken(apiToken).Get(url)

  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  fmt.Println("result", resp.StatusCode())
  fmt.Println("result", resp.String())
  c.Data(resp.StatusCode(), "application/json", []byte(resp.String()))
}

// func deleteFromKv(c *gin.Context){
//   key := c.Param("key")
//   accountId := os.Getenv("KV_ACCOUNT_ID")
//   url := fmt.Sprintf(connectionUrl, accountId, namespaceID, key)

//   resp, err := client.R().SetAuthToken(apiToken).DELETE(url)

//   if err != nil {
//     c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//     return
//   }

//   c.Data(resp.StatusCode(), "application/json", []byte(resp.String()))
// }


func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default environment variables")
	}

  app := gin.Default()

  app.GET("/:key", readFromKv)

  app.POST("/:key", writeToKv)

  // app.DELETE("/:key", deleteFromKv)

  PORT := os.Getenv("PORT")
  if PORT == ""{
    PORT = 
    "8080"
  }
  fmt.Println("Server running on PORT -> ", PORT)
  log.Fatal(app.Run(":" + PORT))
}
