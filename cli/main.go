//go:build linux || windows || amd64

package main

import (
	"embed"
	"goph_keeper/cli/application"
	"goph_keeper/cli/cmd"
	"goph_keeper/cli/storage"
)

//go:embed "migrations/*.sql"
var migrationsFS embed.FS

//go:embed config.json
var configFile []byte

var (
	Version string = "dev"
	Date    string = "now"
)

func main() {
	//TODO: build tags  - cfg + ручка показа билд тегов
	cmd.Version = Version
	cmd.Date = Date
	application.ConfigData = configFile
	storage.MigrationsFS = migrationsFS
	cmd.Execute()
}
