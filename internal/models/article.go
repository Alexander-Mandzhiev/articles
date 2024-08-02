package models

//import "time"

type Article struct {
	ID      string `json:"id,omitempty" db:"id"`
	Title   string `json:"title" db:"title"`
	Text    string `json:"text" db:"text"`
	Authors string `json:"authors" db:"authors"`
	//CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}
