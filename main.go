package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"gopkg.in/yaml.v2"
)

type Banner struct {
	Image string
	Url   string
}

type Site struct {
	settings map[string]string
	Zones    map[string]string
}

type Config struct {
	Defaults map[string]string
	Sites    map[string]Site
	Banners  map[string][]Banner
}

func (c *Config) parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

// TODO: Pre generate Banners
// TODO: Use goroutine

func sample(a []Banner) Banner {
	r := rand.Intn(len(a))
	return a[r]
}

func generateBanner(image string, url string) string {
	return "<a href='" + url + "'><img src='" + image + "'/></a>"
}

func getBanner(site string, zoneID string, c *Config) string {
	b := sample(c.Banners[c.Sites[site].Zones[zoneID]])
	return generateBanner(c.Defaults["base_uri"]+b.Image, b.Url)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	data := `
defaults:
  base_uri: http://static.turfmedia.com/
sites:
  ph:
    settings:
      affiliate-id: 2a848818

    zones:
      top: banner728x90
      middle: banner728x90
banners:
  banner728x90: #skyscraper
    -
      image: turfistar/3481d796.png
      url: "http://turfistar.com/#a_aid=2a848818&a_bid=3481d796"
    -
      image: pronostic-facile/0eca889d.png
      url: "http://www.pronostic-facile.fr?a_aid=2a848818&a_bid=0eca889d"
    -
      image: gazette-turf/8d3376b4.png
      url: "http://gazette-turf.fr/#a_aid=2a848818&a_bid=8d3376b4"
`

	// data, err := ioutil.ReadFile("config.yml")
	// if err != nil {
	// 	println("File does not exist:", err.Error())
	// 	os.Exit(1)
	// }

	// data, err := Asset("config.yml")
	// if err != nil {
	// 	println("config not found:", err.Error())
	// 	os.Exit(1)
	// }

	var config Config
	if err := config.parse([]byte(data)); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", config)

	e := echo.New()
	e.GET("/zones/:site/:id", func(c echo.Context) error {
		zone := c.Param("id")
		site := c.Param("site")
		if config.Sites[site].Zones[zone] != "" {
			return c.String(http.StatusOK, getBanner(site, zone, &config))
		}
		return c.String(http.StatusNotFound, "Not Found")
	})
	e.Static("/", "ads/index.html")
	e.Run(fasthttp.New(":1323"))

}
