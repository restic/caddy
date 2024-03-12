package restic

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	restserver "github.com/restic/rest-server"
)

func init() {
	caddy.RegisterModule(ResticModule{})
}

type ResticModule struct {
	RepositoryPath   string `json:"repository_path,omitempty"`
	AppendOnly       bool   `json:"append_only,omitempty"`
	Debug            bool   `json:"debug,omitempty"`
	MaxRepoSize      int64  `json:"max_repo_size,omitempty"`
	NoVerifyUpload   bool   `json:"no_verify_upload,omitempty"`
	PrivateRepos     bool   `json:"private_repos,omitempty"`
	Prometheus       bool   `json:"prometheus,omitempty"`
	PrometheusNoAuth bool   `json:"prometheus_no_auth,omitempty"`
	HtpasswdPath     string `json:"htpasswd_path,omitempty"`
	NoAuth           bool   `json:"no_auth,omitempty"`

	resticHandler http.Handler
}

func (ResticModule) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.restic",
		New: func() caddy.Module { return new(ResticModule) },
	}
}

func (m *ResticModule) Provision(ctx caddy.Context) error {

	restConfig := restserver.Server{
		Path:             m.RepositoryPath,
		NoAuth:           m.NoAuth,
		HtpasswdPath:     m.HtpasswdPath,
		AppendOnly:       m.AppendOnly,
		Debug:            m.Debug,
		MaxRepoSize:      m.MaxRepoSize,
		NoVerifyUpload:   m.NoVerifyUpload,
		PrivateRepos:     m.PrivateRepos,
		Prometheus:       m.Prometheus,
		PrometheusNoAuth: m.PrometheusNoAuth,
	}

	handler, err := restserver.NewHandler(&restConfig)
	if err != nil {
		return err
	}

	m.resticHandler = handler
	return nil
}

func (m ResticModule) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	m.resticHandler.ServeHTTP(w, r)
	return nil
}

var (
	_ caddy.Provisioner           = (*ResticModule)(nil)
	_ caddyhttp.MiddlewareHandler = (*ResticModule)(nil)
)
