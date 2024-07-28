package models

type TestCreateSchema struct {
	Title string `json:"title"`
}

type TestResponseSchema struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
}
