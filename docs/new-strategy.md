# New strategy:

## Initial download

  - cleanup corrupt folders (from failed attempts....)
  - if not present in vendor/src, download with `go get` - so we dont have to deal with tricky logic here
  - we also generate a lock file after initial download (by traversing the vendor folder for VCS folders)
  - once we have them all (initial download), we copy the .hg/.git/.bzr folders to cache
  - idea: maybe use rsync to for VCS folder sync (this reduces mindless file shuffling)

## Later download
  - update cached VCS folders (git fetch, etc...)
  - rsync to vendor folder
