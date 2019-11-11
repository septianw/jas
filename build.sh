#!/bin/bash

target="."
cwd=$(pwd)

if [[ ! -z $1 ]]
then
  target=$1
fi

targetdir=$target/dist
 
mkdir -p $targetdir
mkdir -p $targetdir/{libs,middlewares,migrations,modules/core,modules/contrib,schema}
cp ./emptyconfig.toml $target/dist/config.toml

go build -ldflags="-s -w" -o $targetdir/jas
upx -9 $targetdir/jas
cp ./schema/* $targetdir/schema/

cd ./migrate && go build -ldflags="-s -w" -o ../$targetdir/migrate && cd ../
upx -9 $targetdir/migrate
