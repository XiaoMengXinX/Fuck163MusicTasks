#!/bin/bash

mkdir out

bash build.sh linux amd64 -o "./out/Fuck163MusicTasks_linux_amd64"
bash build.sh linux arm64 -o "./out/Fuck163MusicTasks_linux_arm64"
bash build.sh windows amd64 -o "./out/Fuck163MusicTasks_windows_amd64.exe"
bash build.sh darwin amd64 -o "./out/Fuck163MusicTasks_darwin_amd64"
bash build.sh darwin arm64 -o "./out/Fuck163MusicTasks_darwin_arm64"
