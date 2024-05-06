package db

type Chapter struct {
	ID         uint   `db:"id"`
	Title      string `db:"title"`
	Content    string `db:"content"`
	RequestURL string `db:"request_url"`
	Number     int    `db:"number"`
	BookId     int    `db:"book_id"`
}

type Novel struct {
	ID     uint   `db:"id"`
	Name   string `db:"name"`
	ImgUrl string `db:"img_url"`
	Img    []byte `db:"img" gorm:"type:bytea"`
	Path   string `db:"path"`
}

func (c Chapter) TableName() string {
	return "chapter"
}

func (c Novel) TableName() string {
	return "novel"
}
