# Gotchas

Forge mounts the user's `docker.sock` into each container that it runs to grant subprocesses access to Docker by default. This feature can be disabled like so:

```sh
forge --no-dind ...
```

However, subprocesses that try to use references to the filesystem when interacting with Docker will run into problems. For example, the following would behave unexpectedly:

```sh
docker run -v /src:/dst
```

This is because the subprocess would likely be using references to the _container's filesystem_ while Docker will interpret them as references to the _host's filesystem_. Forge provides a mechanism that can be enabled to mitigate this shortcoming by intercepting traffic to the mounted `docker.sock` and translating references to mounted directories from the host on the _container's filesystem_ into the _host's filesystem_ equivalent. While this mechanism is not quite stable, it can be enabled as follows:

```sh
forge --use-sock ...
```
