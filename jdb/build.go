package jdb

type Build struct {
	Id     string `json:"id"`
	Passed bool   `json:"passed" binding:"required"`
	Commit string `json:"commit" binding:"required"`
}
