# Hez

# TODO

## Environments

Don't pass environment variables through -e

Mount /etc/hez/environments/[app_name]/.env to /var/www/.env

Make the destination / name of the env file configurable but it just works out of the box by default.

This means that we remove some crap to get the environment => faster, cleaner (stronger?)
You can modify the environment and it instantly updates (but do we want that?)

Maybe implement this system

/etc/hez/environments/someapp/current/.env
/etc/hez/environments/someapp/staging/.env

`hez env commit someapp`

Or

`hez apply --commit-envs`

Or both.

* diagnose command
* ansi output
* documentation
* self-update command / system
* build on release commmand