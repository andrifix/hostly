package path

import (
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"net/http"
	"os"
	"slices"
	"strings"
)

func init() {
	caddy.RegisterModule(&Checker{})
	httpcaddyfile.RegisterHandlerDirective("path_checker", parseCaddyfile)
}

type Checker struct {
	Path    string   `json:"path"`
	Domains []string `json:"domains,omitempty"`
}

func (m *Checker) ServeHTTP(writer http.ResponseWriter, request *http.Request, handler caddyhttp.Handler) error {
	domain := request.URL.Query().Get("domain")
	filename := strings.TrimSuffix(caddyhttp.SanitizedPathJoin(m.Path, domain), "/") + "/"

	if slices.Contains(m.Domains, domain) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("OK"))
		return nil
	}

	if s, err := os.Stat(filename); err == nil && s.IsDir() {
		fmt.Println(s)
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("OK"))
		return nil
	}

	writer.WriteHeader(http.StatusNotFound)
	writer.Write([]byte("NOT FOUND"))
	return nil
}

func (m *Checker) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.path_checker",
		New: func() caddy.Module { return new(Checker) },
	}
}

func (m *Checker) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	d.Next()

	for nesting := d.Nesting(); d.NextBlock(nesting); {
		switch d.Val() {
		case "path":
			if !d.NextArg() {
				return d.ArgErr()
			}
			m.Path = d.Val()
		case "domains":
			m.Domains = d.RemainingArgs()
		default:
			return d.Errf("unrecognized parameter or sub-directive 2 '%s'", d.Val())
		}
	}

	return nil
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Checker
	err := m.UnmarshalCaddyfile(h.Dispenser)

	return &m, err
}
