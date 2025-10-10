package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// Reverse proxy pour l'application Next.js
func ReverseProxy() gin.HandlerFunc {
	target := "localhost:3000"
	url, err := url.Parse("http://" + target)
	if err != nil {
		log.Fatal(err)
	}
	return func(c *gin.Context) {
		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = url.Host
			req.URL.Scheme = url.Scheme
			req.URL.Host = url.Host
			req.URL.Path = c.Request.URL.Path
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	r := gin.Default()

	// Les routes de l'API sont gérées par Go
	api := r.Group("/api")
	{
		api.GET("/message", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Bonjour depuis le backend Go !",
			})
		})
	}

	// Toutes les autres routes sont transmises à l'application Next.js
	// Le proxy doit gérer la racine, les assets, etc.
	r.NoRoute(ReverseProxy())

	r.Run(":8080")
}
