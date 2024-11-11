package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"news/internal/models"
	"news/internal/repository"
	"news/internal/repository/mocks"
	"testing"
)

func TestPosts(t *testing.T) {
	// Создаем mock репозиторий
	mockPost := new(mocks.Post)

	// Настраиваем ожидания для mock репозитория
	mockPost.On("Posts", 2).Return([]models.Post{
		{Title: "Test Post 1", Content: "XXX", PubTime: 1234, Link: "http://test1.com"},
		{Title: "Test Post 2", Content: "ZZZ", PubTime: 4321, Link: "http://test.2com"},
	}, nil)

	// Создаем экземпляр API с mock репозиторием
	api := NewHandler(&repository.Repository{Post: mockPost})

	// Инициализируем роутер
	router := api.InitRouts()

	// Создаем тестовый запрос
	req, err := http.NewRequest(http.MethodGet, "/news/2", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем тестовый ResponseRecorder
	rr := httptest.NewRecorder()

	// Обрабатываем запрос
	router.ServeHTTP(rr, req)

	// Проверяем статус кода
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Проверяем заголовки
	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("handler returned wrong Content-Type: got %v want %v",
			rr.Header().Get("Content-Type"), "application/json")
	}

	// Проверяем тело ответа
	var posts []models.Post
	if err := json.NewDecoder(rr.Body).Decode(&posts); err != nil {
		t.Errorf("handler returned wrong body: %v", err)
	}

	if len(posts) != 2 {
		t.Errorf("handler returned unexpected number of posts: got %v want %v",
			len(posts), 2)
	}

	// Проверяем, что все ожидания были выполнены
	mockPost.AssertExpectations(t)
}
