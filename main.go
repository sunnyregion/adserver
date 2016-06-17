package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Zones map[string][]string
}

func (c *Config) parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func sample(a []string) string {
	r := rand.Intn(len(a))
	return a[r]
}

func main() {
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		println("File does not exist:", err.Error())
		os.Exit(1)
	}

	var config Config
	if err := config.parse(data); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	// gin.SetMode(gin.ReleaseMode)
	r.GET("/zones/:id", func(c *gin.Context) {
		zone := c.Param("id")
		if config.Zones[zone] != nil {
			c.String(http.StatusOK, sample(config.Zones[zone]))
		} else {
			c.String(http.StatusNotFound, "Not Found")
		}
	})
	r.Run() // listen and server on 0.0.0.0:8080
}
