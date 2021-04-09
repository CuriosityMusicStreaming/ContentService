package mysql

import (
	"embed"
	"net/http"
)

// go:embed migrations
var migrations embed.FS

type embedder func()

var MigrationsEmbedder embedder

func (m embedder) GetDir() http.FileSystem {
	return http.FS(migrations)
}
