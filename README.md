# Hez

## Installation

```bash
curl -L https://github.com/redwebcreation/hez/releases/latest/download/hez -o hez
chmod +x ./hez
mv ./hez /usr/bin/hez
```

## Getting started

Hez stores its configuration in `/etc/hez/hez.yml`, you can generate it automatically using the command below:

```bash
hez init
```

This command creates the `/etc/hez/hez.yml` contains the main configuration.

The config file looks like this :

````yaml
applications: [ ]
proxy:
  http:
    enabled: true
    port: 80
  https:
    enabled: true
    port: 443
    self_signed: false
  logs:
    level: 0
    redirections:
      - stdout
````

## Applications

A complex application configuration could look like that :

```yaml
applications:
  - image: example-app
    host: example.com
    container_port: 80
    network: some-network
    warm: true
    volumes:
      - /data/example_com/storage:/var/www/html/storage
    registry:
      username: someone
      password: cool
      server: rg.fr-par.scw.cloud/namespace/app
    env:
      - /dat
      - ...
    hooks:
      container_deployed:
        - php artisan migrate -f
```

Let's break it down, line by line.

The `image` is the docker image of your application, you can also specify a version :

```yaml
example-app:4.2.0
```

The `host` tells the proxy to forward any request for this host to the container on the `container_port`

The `env` key lets you provide environment variables to the image, you can provide one per line using the syntax you are
used to.

```yaml
applications:
  - env:
      - MY_VARIABLE=yes
      - ANOTHER_ONE=true
      - YES="not at all"
```

You can also reference an env file :

```yaml
applications:
  - env:
      - /data/envs/example.env
```

You can now apply your configuration :

```bash
hez apply
```

This command will create all the containers as defined in your configuration. Every time you change your configuration,
you may rerun `hez apply`to apply it.

You may stop all running containers.

```bash
hez stop
```

## Proxy

Hez has an integrated reverse proxy that forwards any request to the right container.

```yaml
proxy:
  http:
    enabled: true
    port: 80
  https:
    enabled: true
    port: 443
  logs: ...
```

You can start it like that :

```bash
hez proxy run
```

You can also register the proxy to run automatically and restart on reboot using systemd:

```bash
hez proxy enable
```

You may disable the systemd integration like so :

```bash
hez proxy disable
```

You may check the status of the proxy by running the following :

```bash
hez proxy status
```

### Logs

The proxy logs received requests. You can configure what it logs and how via your configuration

```yaml
proxy:
  ...
  logs:
    level: 0
    redirections:
      - /tmp/the-proxy.log
      - stdout
```

The log `level` defines the minimal level for a log to get logged. It goes from `-1` to `5`.

So if the level is `4`, logs with a level strictly lower than `4` won't be logged.

Here's a table with the number and their corresponding label :

* `-1` Debug
* `0`  Info
* `1`  Warn
* `2`  Error
* `3`  DPanic
* `4`  Panic
* `5`  Fatal

For the logs that have the minimum required level, you can redirect them to various outputs.

* `stdout` redirects the log to the standard output
* `stderr` redirects the log to the standard error output
* `an absolute path to a file`

appends the log to a file, creates the file if it does not exist.

If you leave `redirections` empty, logs won't be saved.

---

TODO:

* customizable 404, 500 page
* logs redirections validation
* variable interpolation for env: and env files
* if container_port is not open on container x throw an error.
* add request_uri to proxy logs
* websockets, gRPC, HTTP2 (3?)
* container logs?
* check if dns points to the server automatically
* diagnose command

