package types

// Posts Object
type Post struct {
	UserId   int        `json:"userId"`
	Id       int        `json:"id"`
	Title    string     `json:"title"`
	Body     string     `json:"body"`
	Comments []*Comment `json:"comments"`
}

// Comments Object
type Comment struct {
	PostId int    `json:"postId"`
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Body   string `json:"body"`
	Email  string `json:"email"`
}
