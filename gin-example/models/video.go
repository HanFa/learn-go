package models

type Video struct {
	Id          int    `uri:"id"`
	Title       string `form:"title"`
	Description string `form:"description"`
}
