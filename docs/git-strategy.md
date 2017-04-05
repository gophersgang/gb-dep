## Possible ways to deal with GIT repositories

- symlinking? - does not work... GIT ignores symlinks
- https://community.atlassian.com/t5/Git-questions/Can-the-git-folder-be-outwith-of-the-repository/qaq-p/208738


# initial clone
mkdir /tmp/gitplay/cache
git clone https://github.com/immortal/immortal --bare /tmp/gitplay/cache/immortal.git
mkdir -p /tmp/gitplay/src/immortal
ln -s /tmp/gitplay/cache/immortal.git /tmp/gitplay/src/immortal/.git
cd /tmp/gitplay/src/immortal
env GIT_WORK_TREE=/tmp/gitplay/src/immortal git reset --hard 0.10.0
rm /tmp/gitplay/src/immortal/.git





# update
export GIT_WORK_TREE=/tmp/gitplay/src/immortal
export GIT_DIR=/tmp/gitplay/cache/immortal.git
cd /tmp/gitplay/cache/immortal.git && git fetch
cd /tmp/gitplay/src/immortal
git pull
git reset --hard <sha>
