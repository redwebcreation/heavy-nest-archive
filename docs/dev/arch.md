# Architecture

An `Application` has two `ShallowContainer` that contain the `Application` itself + the name of the container.

The name of the container is formed with the host and the container port. If the container is temporary, meaning that it acts as bridge between the current (outdated) active container and the next active container, its name is suffixed with `_temporary`.

An example:
```
Host: vpn.example.com
Port: 8080
Temporary: true
```

Is converted to `vpn_example_com_8080_temporary`.

The name is the only thing used to identify a container and, when applying for the first a configuration with a new application, you must ensure that no other containers are running with the future name of the container's applicaton.

`hez diagnose` can tell you that.

**It is REALLY recommended to run a diagnosis before applying a new configuration**.