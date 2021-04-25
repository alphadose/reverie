#!/bin/bash 

set -e
project_dir=$(pwd)
if [ -f $project_dir/bin/air ]; then
    exit 0
fi

printf "🔨 Installing air\n" 

curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

printf "👍 Done\n"
