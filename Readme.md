# GB Dep

Reliable, fast dependency management for [GB](getgb.io)

## Features
  - `hjson.org` (Human JSON) as package format
  - package.hjson file (a very simple version of the Node.js package.json)
  - supports parallel installation of packages
  - installs binaries in vendor/bin folder

## Non-Features
  - support Golang < 1.8
  - automatic resolution for possible major / minor / patch requirements (you do this by manually)


## Why?
I want to manage my projects with GB, but the GB vendoring tool works only with GIT SHAs. All other verndoring libraries in Golang do too much, are slowish, have lots of features / bugs. I want something very simple.


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

    $ gb-dep install

Installs packages and creates a package.lock file. Next time you run this command, it will take the versions from the package.lock file.

package.lock is tied to the the package.json via a checksum value, that forces package.lock update on package.hjson changes

    $ gb-dep update

forces package.lock update, this will update the not-fixed packages to latest sha on master



## Thx
  - https://github.com/mattn/gom - for some code snippets
  - https://github.com/apoydence/loggr - for subcommands handling logic

## License

MIT licensed. See the LICENSE file for details.