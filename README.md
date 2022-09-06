# forge

## thoughts

the world is made up of CLIs. they're very powerful. they do a lot with a little. we want to use them a lot.

we don't want to use different versions of them, so we want to containerize them

because they do a lot, sometimes they are slow, so we want to cache them

sometimes a container does not exist with all of the tools we need for a sequence of commands, so we want to chain commands+containers together

we want a system to chain them together, but as of yet no such system has proven successful (in my opinion). past attempts include Jenkins, GitHub Actions, etc. they struggle with a singular problem: CLIs are hard to interact with.

CLIs take any number of only string inputs that are then parsed into anything plus a stream of bytes on stdin that are then--again--parsed into anything.

CLIs output their process state, 2 byte streams in the form of stdout and stderr, and any changes they make to the filesystem. sometimes those streams are intended to be parsed into something, others not.

past systems have attempted to remedy this with native wrappers around CLIs e.g. Jenkins Plugins, GitHub Actions (e.g. actions/checkout@v2) (js), Concourse Resources (e.g. registry-image-resource) (go), but that still leaves a lot to be desired as it is inefficient to write native wrappers for any possible CLI someone wants to interact with

it should be possible to run containerized commands, capture their process state, their stdout and stderr streams, and a tar of a specific directory ($OUT a la nix)

all of these should be cacheable via the digest of the image that the command runs in and the digest of the command's encodable representation (e.g `'{ "entrypoint": [...], "cmd": [...] }'`)

a shim that executes the commands inside of the container should be able to do extra cool stuff like run github actions by cloning them, capturing their metadata, and then executing the action (none of this inherently needs to be done by the orchestrator like I did with sequence; moments ago i thought this was important but now I'm questioning it)

concourse resources should be attainable similarly (minus the shim because they are containers in and of themselves)

## try to draw it out

we start with an encodable workload to run:

```json
{
    "image": "sha256:7580ece7963bfa863801466c0a488f11c86f85d9988051a9f9c68cb27f6b7872",
    "cmd": [
        "echo",
        "'hello there'",
        "$WHO"
    ],
    "env": [
        "WHO=GENERALKENOBI"
    ],
    "stdin": "bytes"
}
```

where `image` actually holds an `io.Reader` to get the tar contents of the image from

before running should get the digest of the thing and check to see if it is cached

this should result in the exit code, stdout, stderr, and a tarball of the working directory (or `$OUT` ?)
