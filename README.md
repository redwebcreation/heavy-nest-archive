# Hez

```bash
apt install hez
```

## Getting started

Hez stores all its configuration in `/etc/hez`. It needs quite a few files and directories to work, so you can generate
them automatically using the command below:

```bash
hez config new
```

This command creates the following files :

* `/etc/hez/environments/` contains your applications' environment.
* `/etc/hez/ssl` contains self-generated certificates (if you use them.)
* `/etc/hez/hez.yml` contains the main configuration

You can delete your whole configuration using `hez config delete`.

> This will also stop any running containers and/or the reverse proxy.

The default config file (in `/etc/hez/hez.yml`) looks like this :

````yaml
applications: [ ]
proxy:
  port: 80
  ssl: 443
  self_signed: false
logs:
  level: 0
  redirections:
    - for: out
      value: stdout
    - for: err
      value: stderr
````

Let's break it down.

## Logs

The log `level` defines the minimal level for a log to get logged. It goes from `-1` to `5`.

So if the level is `4`, logs with a level strictly lower than `4` won't be logged.

Here's a table with the number and their corresponding label :

* `-1`: Debug
* `0`: Info
* `1`: Warn
* `2`: Error
* `3`: DPanic
* `4` Panic
* `5` Fatal

Now, for the logs that have the minimum required level, you can redirect them to various outputs.

There's two type of logs, logs coming from the standard output, and those coming for the standard error output.

You may specify which one you want with the `for` key which accepts either `out` or `err`

Now, to redirect the type of log you chosen, you can specify a `value` which can be to the following :

* `stdout` redirects the log to the standard output
* `stderr` redirects the log to the standard error output
* `an absolute path to a file` appends the log to a file, creates the file if it does not exist.

If you leave `redirections` empty, no log will be saved.

## Applications

You can create a new application by adding an entry to the `applications:` section of your config.

```yaml
applications:
  - name: yourAppName
    image: nginx
    environment: nginx.env
    bindings:
      - host: example.com
        port: 8000
```

Let's break this down, line by line:

The `applications` key just contains the list of your applications, pretty straight-forward.

Now, the name, it will be removed soon, but right now it's used to show the app's name. It can be whatever, it's not
actually used.

The image represents the docker image of your application, you can specify a version like this:

```yaml
applications:
  - ...
    image: nginx:latest
```

Just as you would do with Docker.

The environment will get completely reworked tomorrow (probably), so no docs for now.

The bindings are probably the most complicated part of this configuration :

You bind `http://example.com` and `https://example.com` to your container's port 8000. You're not
binding `http://example.com:8000` to your container on a port.

You can apply your configuration :

```bash
hez apply
```

If you didn't change your config and try to apply, it won't work.

```bash
hez apply -f
```

You need to force it.

You can also stop running containers

```bash
hez stop
```

**TODO: Prevent stopping if proxy is running or option --with-proxy**

You probably want to stop your proxy before.

## Proxy

There's an integrated reverse proxy with sub-millisecond response time integrated. (It takes a ~1ms with 100 containers
and the worst case scenario)

You can start it like that

```bash
hez proxy run
```

In development, you probably want to do smtg like

```bash
hez proxy run --port 8080 --ssl 8443 --self-signed
```

It generates SSL certificates automatically (and re-generates) using Let's Encrypt.

It does that one the first request, and then everytime the certificates expires it regenerates it during a request.

// TODO: A registrable cron command or smtg

If you pass no options `hex proxy run`, it uses the options in your configuration `proxy:`.

You can register / disable the reverse proxy in systemctl. It can run / restart automatically.

```bash
hez proxy enable
```

```bash
hez proxy disable
```

You can also get the proxy's status with `hez proxy status`.

The proxy is registered as `hezproxy.service` in systemd.

# TODO

## Environments

Don't pass environment variables through -e

Mount /etc/hez/environments/[app_name]/.env to /var/www/.env

Make the destination / name of the env file configurable but it just works out of the box by default.

This means that we remove some crap to get the environment => faster, cleaner (stronger?)
You can modify the environment and it instantly updates (but do we want that?)

Maybe implement this system

/etc/hez/environments/someapp/current/.env /etc/hez/environments/someapp/staging/.env

`hez env commit someapp`

Or

`hez apply --commit-envs`

Or both.

* GetConfig returns a Config instead of Config, error
* diagnose command
* ansi output
* documentation
* self-update command / system
* build on release commmand