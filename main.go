package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"gopkg.in/yaml.v2"
)

type Banner struct {
	Image string
	Url   string
}

type Config struct {
	Zones   map[string][]int
	Banners map[int]Banner
}

func (c *Config) parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

// TODO: Pre generate Banners
// TODO: Use goroutine

func sample(a []int) int {
	r := rand.Intn(len(a))
	return a[r]
}

func generateBanner(image string, url string) string {
	return "<a href='" + url + "'><img src='" + image + "'/></a>"
}

func getBanner(zoneID string, c *Config) string {
	b := c.Banners[sample(c.Zones[zoneID])]
	return generateBanner(b.Image, b.Url)
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

	fmt.Printf("%#v\n", config)

	e := echo.New()
	e.GET("/zones/:id", func(c echo.Context) error {
		zone := c.Param("id")
		if config.Zones[zone] != nil {
			return c.String(http.StatusOK, getBanner(zone, &config))
		}
		return c.String(http.StatusNotFound, "Not Found")
	})
	e.Static("/", "public/index.html")
	e.Run(fasthttp.New(":1323"))

}
