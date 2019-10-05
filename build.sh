#!/bin/bash

PACKAGE_NAME='gospeckle'

if [ -n "$1" ]
then
  VERSION=$1
else
  echo "A version number must be supplied"
  exit 1
fi

GIT_COMMIT=$(git rev-list -1 HEAD)
# DATE=`date '+%d-%m-%Y'`
DATE=`date -u`

platforms=("windows/amd64" "windows/386" "darwin/amd64" "linux/amd64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    output_name='dist/'$PACKAGE_NAME'-'$GOOS'-'$GOARCH

    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH \
      go build \
      -ldflags "-X github.com/speckleworks/gospeckle/cmd.Version=$VERSION \
                -X github.com/speckleworks/gospeckle/cmd.GitCommit=$GIT_COMMIT \
                -X github.com/speckleworks/gospeckle/cmd.GitCommit=$GIT_COMMIT \
                -X \"github.com/speckleworks/gospeckle/cmd.BuildDate=$DATE\" \
                -X github.com/speckleworks/gospeckle/cmd.OsArch=$platform" \
      -o $output_name
done