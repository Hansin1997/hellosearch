package main

import (
	"context"
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
)

type Config struct {
	Host string `json:"host"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

func loadConfig(file string) (*Config, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	c := Config{}
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	storage, err := NewElasticStorage([]string{config.Host},
		config.User,
		config.Pass,
		context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/doc", func(c *gin.Context) {
		q, e := url.QueryUnescape(c.Query("q"))
		if e != nil {
			q = c.Query("q")
		}
		p, e := strconv.Atoi(c.Query("p"))
		if e != nil {
			p = 1
		}
		r, e := storage.Search(q, c.QueryArray("i"), p)
		if e != nil {
			c.JSON(400, e)
		} else {
			c.JSON(200, r)
		}
	})
	r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
