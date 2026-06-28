package public

import "embed"

//go:embed all:*
var StaticFS embed.FS

//go:embed xevon-configs.example.yaml
var DefaultConfigYAML []byte

//go:embed workbench-users.json
var WorkbenchUsersJSON []byte
