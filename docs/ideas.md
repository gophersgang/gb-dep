### Ideas to improve dependency resolution

  - take hints from the DB world and introduce a concept of a `query plan`
  - from devops: converge on a desired state. Have a proper dependency graph
  - testing: introduce a mock backend for deps, that would allow testing complex scenarios, with time-forwarding, dependency changes and such
  - debugging: show the planned execution plan in `dry` run



  - introspection:
    - show dependency graph
    - show size of dependencies
    -


### Vendoring:
  - helps dealing with a `vendor` - repo in a consistent / VCS-independent way
  - provides configuration for vendoring, like
    - vendor-repo: github.com/user/project-vendor.git
    - vendor-commit: reference to a commit

  - store exact metadata in the vendor repo
  - that allows exact comparison with the lock file in the main project repo



### Commands to deal with the vendor repo:
    $ gbdep vendor init
      - will cleanup vendor/src from .git / .hg / .bzr folders
      - copy `package.lock` to `vendor/src` (to know, wherether exact match is given)
      - vendor/src becomes a git repo
      -

    $ gbdep vendor fetch
