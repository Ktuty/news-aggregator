package models

type Post struct {
	Id      int    //`json:"id"`
	Title   string //`json:"title"`
	Content string //`json:"content"`
	PubTime int64  //`json:"publication_date"`
	Link    string //`json:"link"`
}
