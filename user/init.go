package user

import "github.com/jeyem/mogo"

var (
	db *mogo.DB
)

func Register(database *mogo.DB) {
	db = database

}
