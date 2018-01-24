package file

import "github.com/jeyem/mogo"

var (
	db     *mogo.DB
	config *Config
)

type Config struct {
	ImagePath     string
	ServingPrefix string
}

func Register(database *mogo.DB, conf Config) {
	db = database
	config = &conf

}
