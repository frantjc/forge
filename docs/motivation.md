# Why?

What problem does this solve? Well...

Automation begins with a shell script that executes a bunch of CLI commands often to test, build and publish some code. The next step is to set up some continuous integration (CI) system that executes that script in response to some event such as a commit to a Git repository's `main` branch. Such CI systems tend to identify that all of the scripts that they are executing do a lot of the same things--checkout a Git repository, setup a tool and so on.

In an effort to make their platform easier to use and to refactor the shared functionality out of all of the aforementioned scripts, CI systems in the past have introduced reusable "Actions"/"CloudBuilders"/"resources"/"tasks" which take minimal configuration to do a complex task. GitHub Actions' [`actions/checkout`](https://github.com/actions/checkout), for example, takes one short line of code to invoke and accepts a bunch of optional configuration to fulfill many related use cases.

Unfortunately, using such powerful plugins outside of the the system they were built for can be wildly difficult. This makes debugging the use of these plugins require long feedback loops. It also makes migrating from one CI system to another treacherous, having to replace uses of one system's plugins with another's.

Forge aims to remedy this.
