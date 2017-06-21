package restic

import (
	"net/http"
	"strings"

	"github.com/mholt/caddy/caddyhttp/httpserver"
)

type ResticHandler struct {
	Next          httpserver.Handler
	BasePath      string
	RestServerMux http.Handler
}

func (h ResticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if !httpserver.Path(r.URL.Path).Matches(h.BasePath) {
		return h.Next.ServeHTTP(w, r)
	}

	// basic auth is required (some authentication is required, for obvious reasons)
	basicAuthUser, ok := r.Context().Value(httpserver.RemoteUserCtxKey).(string)
	if !ok || basicAuthUser == "" {
		return http.StatusForbidden, nil
	}

	// strip the base path from the request so that the restic mux can load
	// the repo relative to its Config.Path base path.
	r.URL.Path = strings.TrimPrefix(r.URL.Path, h.BasePath)
	if !strings.HasPrefix(r.URL.Path, "/") {
		r.URL.Path = "/" + r.URL.Path
	}

	// TODO: this doesn't return values so errors may not be handled properly. Oh well.
	h.RestServerMux.ServeHTTP(w, r)

	return 0, nil
}
