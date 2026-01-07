package domain

type Comment struct {
	BaseModel

	UserID  uint64
	PostID  uint64
	Comment string
	Rating  *uint8

	User *User
}
