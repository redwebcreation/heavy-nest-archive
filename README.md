# Hez

**This documentation is awful, it's just me putting together some sentences from the future documentation.**

**Even worse, I'm pretty sure it's not English down there, don't read it. It's just bad.**

## Installation

// Installation instructions here

```bash
hez config new
```

Generates all needed files in `/etc/hez` with the right permissions (this can be tricky to do yourself)

You can also remove all  the generated config files with 

```bash
hez config delete
```

**RIGHT NOW THE STORAGE DIRECTORY IN $HOME/.config/hez/storage will not be deleted. It's a bug.**

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

## Logs

You can configure logs in your config file

The level goes from -1 (debug) to 8 (or 7 need to check) Fatal.

// List of all levels and their associated label like (-1: debug, 0: info...)

You can redirect the output and the error output to a file or you can redirect it to stdout, stderr

```yaml
redirections:
  - for: out
    value: stdout
```

Redirects stdout to stdout, by default if `redirections` is empty, nothing is redirected / saved. Nothing happens.

```yaml
redirections:
  - for: out
    value: /tmp/my-log-file.log
```

Saves all the logs in a file. Useful for things like Kibana...

You can also use `for: err` which redirects everything coming in stderr.

That's it and its already pretty cool.


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