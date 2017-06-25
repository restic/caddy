Caddy restic plugin
===================

This plugin makes it easy to run a [restic backup](https://github.com/restic/restic) server! This plugin uses [restic/rest-server](https://github.com/restic/rest-server) to make your backup repositories reachable over HTTPS.

Using restic's "rest" backend instead of the "sftp" backend is likely to provide faster transfer speeds because it avoids a lot of SFTP's flow control problems, where transfers continually get slower and slower forever.

The advantage of using this plugin over the bare `rest-server` command is that Caddy provides HTTPS by managing TLS certificates for you, so you always get a secure access point for your repositories.

## Usage

The underlying rest-server needs a data path within which the restic repositories are stored. By default, it is `/tmp/restic`, but you should change this to a non-temporary path for storing actual backups.

**The `restic` plugin _requires_ authentication using [basicauth](https://caddyserver.com/docs/basicauth)** because otherwise anyone could access your backups. Yes, restic backups are encrypted, but people could still delete them, etc. The examples here are shown in conjunction with the `basicauth` directive for this reason. Requests to the rest-server endpoints without authentication will be forbidden. **Caddy uses HTTPS by default, and it is not safe to use this plugin without HTTPS (TLS). Do not disable TLS.**

The syntax of the `restic` directive is:

```
restic [base_path [data_path]]
```

The `restic` plugin can be used without arguments:

```
basicauth / user pass
restic
```

This configures all requests to be processed by rest-server using the default data path. This is useful mostly for testing or experimenting.

```
basicauth /backups user pass
restic    /backups
```

This does the same thing, but only for HTTP requests to `/backups`. The data path is the same.

```
basicauth /backups user pass
restic    /backups /home/me/backups
```

This sets the base path of requests to be `/backups` and the data path on disk to be `/home/me/backups`. In other words, all requests to `/backup` will be handled by the rest-server to manage restic repositories within `/home/me/backups`.

Once your server is running, you can access your backups via HTTPS with restic quite easily:

```
$ restic --repo "rest:https://user:pass@example.com/repo_name" snapshots
```

The path to the repository (`repo_name`) will be relative to the data path you specified in your Caddyfile.
