package main

import (
	"github/Panyakorn4/kwanjai-shop-tutorial/config"
	"github/Panyakorn4/kwanjai-shop-tutorial/modules/servers"
	"github/Panyakorn4/kwanjai-shop-tutorial/pkg/databases"
	"os"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := databases.DbConnect(cfg.Db())
	defer db.Close()

	servers.NewServer(cfg, db).Start()
}
