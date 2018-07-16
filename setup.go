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

	restConfig := restserver.Server{}

	for c.Next() {
		if c.NextArg() {
			basePath = c.Val()
		}
		if c.NextArg() {
			restConfig.Path = c.Val()
		}
		if c.NextArg() {
			return c.ArgErr()
		}
	}

	if basePath == "" {
		basePath = "/"
	}

	httpConfig := httpserver.GetConfig(c)
	mid := func(next httpserver.Handler) httpserver.Handler {
		return ResticHandler{
			Next:          next,
			BasePath:      basePath,
			RestServerMux: restserver.NewHandler(restConfig),
		}
	}
	httpConfig.AddMiddleware(mid)

	return nil
}
