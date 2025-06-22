package res

import (
	"embed"
)

// Content holds the embedded migrations.
//
//go:embed migrations/*.sql
var Content embed.FS
