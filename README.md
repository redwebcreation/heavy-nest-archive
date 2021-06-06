# Hez

```bash
apt install hez
```

#### Disclaimer

Some commands might seem familiar with what [Kubernetes](https://kubernetes.io) does. It's a pure coincidence so don't
expect things to work the same as Kubernetes.

## Getting started

Hez stores its configuration in `/etc/hez/hez.yml`, you can generate it automatically using the command below:

```bash
hez config new
```

This command creates the `/etc/hez/hez.yml` contains the main configuration.

You may delete your configuration file using `hez config delete`.

> This will also stop any running containers and/or the reverse proxy.

The config file looks like this :

````yaml
applications: [ ]
proxy:
  port: 80
  ssl: 443
  self_signed: false
  logs:
    level: 0
    redirections:
      - stdout
````

## Applications

An application in the configuration looks like that :

```yaml
applications:
  - image: example-app
    host: example.com
    container_port: 80
    env:
      - APP_ENV=local
      - ...
```

Let's break it down, line by line.

The `image` is the docker image of your application, you can also specify a version :

```yaml
example-app:4.2.0
```

The `env` key lets you provide environment variables to the image, you can provide one per line using the syntax you are
used to.

```yaml
applications:
  - env:
      - MY_VARIABLE=yes
      - ANOTHER_ONE=true
      - YES="not at all"
```

The `host` tells the proxy to forward any request for this host to the application on the `container_port`


## Proxy

Hez has an integrated reverse proxy that forwards any request to the right container.

You can start it like that :

```bash
hez proxy run
```

You can also specify the ports that the proxy should listen to.

```bash
hez proxy run --port 8080 --ssl 8443
```

By default, the proxy will use Let's Encrypt to generate (and re-generate) SSL certificates, but you may want to use
self-signed certificates for testing :

````bash
hez proxy run --self-signed
````

> You can not use Let's Encrypt to secure localhost.


By default the ports used, and the SSL strategy is defined in your configuration file :

```yaml
proxy:
  port: 80
  ssl: 443
  self_signed: true
  logs: ...
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
* `an absolute path to a file` appends the log to a file, creates the file if it does not exist.

If you leave `redirections` empty, logs won't be saved.

---
**It's just notes from now on. It's outdated and contains a lot of tpyos.**

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

TODO:

* check if dns points to the server automatically
* diagnose command
* ansi output
* documentation
* self-update command / system
* build on release commmand
