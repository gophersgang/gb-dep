# GB Dep

Reliable, fast dependency management for [GB](https://getgb.io/)

## Features
  - `hjson.org` (Human JSON) as package format
  - package.hjson file (a very simple version of the Node.js package.json)
  - supports parallel installation of packages
  - installs binaries in vendor/bin folder

## Non-Features
  - support Golang < 1.8
  - automatic resolution for possible major / minor / patch requirements (there is not good way for this right now)


## Why?
I want to manage my projects with GB, but the GB vendoring tool works only with GIT SHAs. All other verndoring libraries in Golang do too much, are slowish, have lots of features / bugs. I want something very simple.

## Installation


    $ go get -u github.com/gophersgang/gbdep/...


## Usage
Create a package.hjson file with your packages like:

```
packages: [
  // default packages
  { name: "github.com/gorilla/mux", tag: "v1.3.0" }

  // dev packages
  { name: "github.com/mattn/gover", commit: "x8948594854" , group: ["development"], goos: [ "windows", "linux", "darwin" ] }

  // test packages
  { name: "github.com/mattn/gom", commit: "x8948594854", group: ["test"], goos: [ "windows", "linux", "darwin" ] }
]
```

    $ gbdep install

Installs packages and creates a package.lock file. Next time you run this command, it will take the versions from the package.lock file.

package.lock is tied to the the package.json via a checksum value, that forces package.lock update on package.hjson changes

    $ gbdep update

forces package.lock update, this will update the not-fixed packages to latest sha on master


    $ gbdep buildbins

In case you have deleted vendor/bin folder, this will recompile all the binaries. The install command does it only on first pass, to not slow down the operation without much gain.


## Thx
  - https://github.com/mattn/gom - for some code snippets
  - https://github.com/apoydence/loggr - for subcommands handling logic

## License

MIT licensed. See the LICENSE file for details.