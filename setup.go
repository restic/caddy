package restic

import (
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	restserver "github.com/restic/rest-server"
)

func init() {
	caddy.RegisterPlugin("restic", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	var basePath string

	for c.Next() {
		if c.NextArg() {
			basePath = c.Val()
		}
		if c.NextArg() {
			restserver.Config.Path = c.Val()
		}
		if c.NextArg() {
			return c.ArgErr()
		}
	}

	if basePath == "" {
		basePath = "/"
	}

	cfg := httpserver.GetConfig(c)
	mid := func(next httpserver.Handler) httpserver.Handler {
		return ResticHandler{
			Next:          next,
			BasePath:      basePath,
			RestServerMux: restserver.NewMux(),
		}
	}
	cfg.AddMiddleware(mid)

	return nil
}
