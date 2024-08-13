Caddy restic plugin
===================

This plugin makes it easy to run a [restic backup](https://github.com/restic/restic) server! This plugin uses [restic/rest-server](https://github.com/restic/rest-server) to make your backup repositories reachable over HTTPS.

Using restic's "rest" backend instead of the "sftp" backend is likely to provide faster transfer speeds because it avoids a lot of SFTP's flow control problems, where transfers slow down more than necessary.

The advantage of using this plugin over the bare [restic/rest-server](https://github.com/restic/rest-server) command is that Caddy provides HTTPS by managing TLS certificates for you, so you always get a secure access point for your repositories and you don't have to reload the server to renew certificates.

## Configuration
The native configuration approach for caddy has changed from caddyfile to json. Even though the caddyfile is still supported the preferred configuration approach is json based. Hence, just the json based configuration is documented here.

The restic plugin is implemented as middleware handler and can be plugged into the middleware pipeline of caddy via configuration. A simple sample configuration is shown below by adding the restic handler to the pipeline:

```
{
  "apps": {
    "http": {
      "servers": {
        "restic": {
          "routes": [{
            "handle": [{
              "handler": "restic",
              "repository_path": "path to the repository"
            }]
          }]
        }
      }
    }
  }
}
```

All significant parameters defined by the [restic/rest-server](https://github.com/restic/rest-server) are available for this plugin too. The following parameters can be defined:

* repository_path: data directory (default "/tmp/restic")
* append_only: enable append only mode
* debug: output debug messages
* max_repo_size: the maximum size of the repository in bytes
* no_verify_upload: do not verify the integrity of uploaded data. DO NOT enable unless the rest-server runs on a very low-power device
* private_repos: users can only access their private repo
* prometheus: enable Prometheus metrics
* prometheus_no_auth: disable auth for Prometheus /metrics endpoint
* htpasswd_path: location of .htpasswd file (default: "\<data directory\>/.htpasswd")
* no_auth: disable .htpasswd authentication

## Authentication
It is highly recommended to require authentication to access the repository. Otherwise anyone could access your backups. (Yes, restic backups are encrypted, but people could still delete them, etc.). Furthermore, since basic authentication is used all communication must be encrypted to protect the credentials sent to the server. **Caddy uses HTTPS by default, and it is not safe to use this plugin without HTTPS (TLS). Do not disable TLS.**

Generally, there are two options to configure authentication:
* The authentication is by default still handled by [restic/rest-server](https://github.com/restic/rest-server) using the .htaccess file which is stored in the base repository directory. Please consult the official [restic/rest-server](https://github.com/restic/rest-server) documentation about how to create such an .htaccess file.
* It is also possible to configure caddy to request the authentication before actually invoking the [restic/rest-server](https://github.com/restic/rest-server). Please consult the official documentation about the configuration possibilities. Once caddy is configured to handle the authentication you can disable the authentication for [restic/rest-server](https://github.com/restic/rest-server) by configuring the no_auth option.

## Access the repository
Once your server is running, you can access your backups via HTTPS with restic quite easily:

```
$ restic --repo "rest:https://user:pass@example.com/repo_name" snapshots
```