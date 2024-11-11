package rss

import (
	"news/internal/repository"
	"news/internal/repository/mocks"
	"os"
	"testing"
	"time"
)

var configPath = "C:/Users/User/GolandProjects/github.com/Ktuty/news-aggregator/internal/rss/sites.json"

func TestNewRSS(t *testing.T) {
	// Создаем mock репозиторий
	mockPost := new(mocks.Post)

	// Создаем временный файл для конфигурации
	configData := []byte(`{
		"rss": [
			"https://habr.com/ru/rss/hub/go/all/?fl=ru",
			"https://habr.com/ru/rss/best/daily/?fl=ru",
			"https://cprss.s3.amazonaws.com/golangweekly.com.xml",
			"https://go.dev/blog/feed.atom?format=xml",
			"https://blog.jetbrains.com/go/feed/"
		],
		"request_period": 1
	}`)
	tmpFile, err := os.CreateTemp("", "sites.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(configData); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	// Переопределяем путь к конфигурационному файлу
	originalPath := configPath
	defer func() {
		configPath = originalPath
	}()
	configPath = tmpFile.Name()

	// Создаем экземпляр RSS с mock репозиторием
	rss := NewRSS(&repository.Repository{Post: mockPost})

	if rss == nil {
		t.Error("NewRSS returned nil")
	}

	if len(rss.links) != 5 {
		t.Errorf("expected 5 links, got %d", len(rss.links))
	}

	if rss.requestPeriod != 1 {
		t.Errorf("expected request period 1, got %d", rss.requestPeriod)
	}
}

func TestStartPolling(t *testing.T) {
	// Создаем mock репозиторий
	mockPost := new(mocks.Post)

	// Создаем экземпляр RSS с mock репозиторием
	rss := &RSS{
		repo:          &repository.Repository{Post: mockPost},
		links:         []string{"http://example.com/rss"},
		requestPeriod: 1,
	}

	// Запускаем StartPolling в отдельной горутине
	go rss.StartPolling()

	// Даем время для выполнения StartPolling
	time.Sleep(2 * time.Second)

	// Проверяем, что все ожидания были выполнены
	mockPost.AssertExpectations(t)
}

func TestCheckLink(t *testing.T) {
	// Создаем mock репозиторий
	mockPost := new(mocks.Post)

	// Создаем экземпляр RSS с mock репозиторием
	rss := &RSS{
		repo:          &repository.Repository{Post: mockPost},
		links:         []string{"http://example.com/rss"},
		requestPeriod: 1,
	}

	// Запускаем CheckLink
	rss.CheckLink()

	// Даем время для выполнения CheckLink
	time.Sleep(1 * time.Second)

	// Проверяем, что все ожидания были выполнены
	mockPost.AssertExpectations(t)
}
