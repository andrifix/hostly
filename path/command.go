package path

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/caddyserver/caddy/v2"
	caddycmd "github.com/caddyserver/caddy/v2/cmd"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
)

//go:embed minimal-config.json
var config embed.FS

func init() {
	caddycmd.RegisterCommand(caddycmd.Command{
		Name:  "multi-file-server",
		Func:  cmdHello,
		Usage: "[--root <path>]",
		Short: "Multi domain file server",
		CobraFunc: func(cmd *cobra.Command) {
			cmd.Flags().StringP("root", "r", "./", "The path to the root of the site")
			cmd.RunE = caddycmd.WrapCommandFuncForCobra(cmdHello)
		},
	})
}

func cmdHello(flags caddycmd.Flags) (int, error) {
	cfg, err := config.ReadFile("minimal-config.json")
	if err != nil {
		return caddy.ExitCodeFailedStartup, fmt.Errorf("cant read config file: %v", err)
	}

	root := flags.String("root")
	if !strings.HasPrefix(root, string(filepath.Separator)) {
		root += string(filepath.Separator)
	}

	rootEscaped, err := json.Marshal(flags.String("root"))
	if err != nil {
		return caddy.ExitCodeFailedStartup, fmt.Errorf("cant marshal root string: %v", err)
	}

	rootEscaped = bytes.Trim(rootEscaped, "\"")

	const PlaceholderPath = "{{base_path}}"

	cfg = bytes.ReplaceAll(cfg, []byte(PlaceholderPath), rootEscaped)

	err = caddy.Load(cfg, true)
	if err != nil {
		return caddy.ExitCodeFailedStartup, err
	}

	select {}
}
