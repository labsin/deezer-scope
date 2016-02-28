#!/bin/bash
pushd po
echo "Generate mo files"
./generate_mo.sh
echo "Update ini translations"
./translate_ini.sh
popd
pushd deezer
echo "building current arch"
go build && mv deezer x86_64-linux-gnu/deezer-scope
echo "build current arch"
echo "building arm"
#click chroot -a armhf -f ubuntu-sdk-15.04 -s vivid run CGO_ENABLED=1 GOARCH=arm GOARM=7 PKG_CONFIG_LIBDIR=/usr/lib/arm-linux-gnueabihf/pkgconfig:/usr/lib/pkgconfig:/usr/share/pkgconfig GOPATH=~/dev-go CC=arm-linux-gnueabihf-gcc CXX=arm-linux-gnueabihf-g++ go build -ldflags '-extld=arm-linux-gnueabihf-g++' && mv deezer arm-linux-gnueabihf/deezer-scope
click chroot -a armhf -f ubuntu-sdk-15.04 -s vivid run CGO_ENABLED=1 GOARCH=arm GOARM=7 GOPATH=$GOPATH PKG_CONFIG_LIBDIR=/usr/lib/arm-linux-gnueabihf/pkgconfig:/usr/lib/pkgconfig:/usr/share/pkgconfig CC=arm-linux-gnueabihf-gcc CXX=arm-linux-gnueabihf-g++ go build -ldflags '-extld=arm-linux-gnueabihf-g++' && mv deezer arm-linux-gnueabihf/deezer-scope
echo "build arm"
popd
pushd ../
echo "building click"
click build deezer-scope
echo "build click"
popd
