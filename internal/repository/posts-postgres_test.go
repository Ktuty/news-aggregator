package repository

import (
	"news/internal/models"
	"news/internal/repository/mocks"
	"reflect"
	"testing"
)

func TestPostPostgres_CreatePost(t *testing.T) {
	type args struct {
		post models.Post
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Create",
			args: args{
				post: models.Post{
					Title:   "Test Title",
					Content: "Test Content",
					PubTime: 123,
					Link:    "http://test.com",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPost := mocks.NewPost(t)
			mockPost.On("CreatePost", tt.args.post).Return()

			r := Repository{
				Post: &PostPostgres{db: nil},
			}

			r.Post = mockPost // Устанавливаем мок

			r.CreatePost(tt.args.post)

			mockPost.AssertExpectations(t)
		})
	}
}

func TestPostPostgres_Posts(t *testing.T) {
	type args struct {
		quantity int
	}
	tests := []struct {
		name    string
		args    args
		want    []models.Post
		wantErr bool
	}{
		{
			name: "Get Posts",
			args: args{
				quantity: 5,
			},
			want: []models.Post{
				{
					Id:      1,
					Title:   "Test Title 1",
					Content: "Test Content 1",
					PubTime: 456,
					Link:    "http://test1.com",
				},
				{
					Id:      2,
					Title:   "Test Title 2",
					Content: "Test Content 2",
					PubTime: 789,
					Link:    "http://test2.com",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPost := mocks.NewPost(t)
			mockPost.On("Posts", tt.args.quantity).Return(tt.want, nil)

			r := Repository{
				Post: &PostPostgres{db: nil},
			}
			r.Post = mockPost // Устанавливаем мок

			got, err := r.Posts(tt.args.quantity)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostPostgres.Posts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostPostgres.Posts() = %v, want %v", got, tt.want)
			}

			mockPost.AssertExpectations(t)
		})
	}
}
