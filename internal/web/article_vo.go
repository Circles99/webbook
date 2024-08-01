package web

type LikeReq struct {
	Like bool  `json:"like"`
	Id   int64 `json:"id"`
}

type ArticleVo struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	Abstract   string `json:"abstract"`
	Content    string `json:"content"`
	AuthorId   int64  `json:"authorId"`
	AuthorName string `json:"authorName"`
	Status     uint8  `json:"status"`
	Created    string `json:"created"`
	Updated    string `json:"updated"`
}

type ListReq struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
