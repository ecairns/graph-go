package app

import (
	"database/sql"
	"fmt"
	"github.com/BurntSushi/toml"
	_ "github.com/lib/pq"
)

type database struct {
	Server   string
	Username string
	Password string
	Database string
	Port     string
}

type TomlConfig struct {
	Database database
}

var Config TomlConfig
var DB *sql.DB

func init() {

	if _, err := toml.DecodeFile("config.toml", &Config); err != nil {
		fmt.Println(err)
		return
	}

}
func DbInit() {
	var err error

	connStr := fmt.Sprintf("postgres://%v:%v@%v/%v",
		Config.Database.Username,
		Config.Database.Password,
		Config.Database.Server,
		Config.Database.Database,
	)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func DbClose() {
	fmt.Println("***** CLOSING DB CONNECTION *****")
	DB.Close()
}
