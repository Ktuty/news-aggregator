package rss

import (
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"log"
	"news/internal/models"
	"news/internal/repository"
	"os"
	"time"
)

type Config struct {
	RSS           []string `json:"rss"`
	RequestPeriod int      `json:"request_period"`
}

type RSS struct {
	repo          *repository.Repository
	links         []string
	requestPeriod int
}

var (
	errorsChannel = make(chan error)
	postsChannel  = make(chan models.Post)
)

// конструктор для создания экземпляра RSS
func NewRSS(repo *repository.Repository) *RSS {
	data, err := os.ReadFile("C:/Users/User/GolandProjects/github.com/Ktuty/news-aggregator/internal/rss/sites.json")
	if err != nil {
		errorsChannel <- fmt.Errorf("error reading config file: %v", err)
		return nil
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		errorsChannel <- fmt.Errorf("error unmarshalling JSON: %v", err)
		return nil
	}

	return &RSS{
		repo:          repo,
		links:         config.RSS,
		requestPeriod: config.RequestPeriod,
	}
}

// Прослушивание каналов
func (rss RSS) StartPolling() {
	ticker := time.NewTicker(time.Duration(rss.requestPeriod) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rss.CheckLink()

		case err := <-errorsChannel:
			log.Printf(err.Error())

		case post := <-postsChannel:
			rss.repo.CreatePost(post)
		}
	}
}

// Запуск обхода лент для извлечения данных из RSS
func (rss RSS) CheckLink() {
	fmt.Print("\n")
	for _, link := range rss.links {
		log.Printf("Checking link: %s", link)
		go rss.parseFeed(link)
	}
}

func (rss RSS) parseFeed(url string) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		errorsChannel <- fmt.Errorf("error parsing feed %s: %v", url, err)
		return
	}

	for _, item := range feed.Items {
		post := models.Post{
			Title:   item.Title,
			Content: item.Content,
			PubTime: item.PublishedParsed.Unix(),
			Link:    item.Link,
		}
		postsChannel <- post
	}
}
