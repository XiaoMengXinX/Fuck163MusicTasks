#!/bin/bash

VERSION=$(git describe --tags)

mkdir out

bash build.sh linux amd64 -o "./out/Fuck163MusicTasks_${VERSION}_linux_amd64"
bash build.sh linux arm64 -o "./out/Fuck163MusicTasks_${VERSION}_linux_arm64"
bash build.sh windows amd64 -o "./out/Fuck163MusicTasks_${VERSION}_windows_amd64.exe"
bash build.sh darwin amd64 -o "./out/Fuck163MusicTasks_${VERSION}_darwin_amd64"
bash build.sh darwin arm64 -o "./out/Fuck163MusicTasks_${VERSION}_darwin_arm64"
