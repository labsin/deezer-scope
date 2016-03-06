#!/bin/bash
click chroot -a armhf -f ubuntu-sdk-15.04 -s vivid maint apt-get update
click chroot -a armhf -f ubuntu-sdk-15.04 -s vivid maint apt-get install -y golang-go golang-go-linux-arm libglib2.0-dev:armhf crossbuild-essential-armhf bzr libaccounts-glib-dev:armhf
mkdir temp
pushd temp
click chroot -a armhf -f ubuntu-sdk-15.04 -s vivid maint apt-get download libsignon-glib-dev:armhf libsignon-glib1:armhf signond-dev:armhf
click chroot -a armhf -f ubuntu-sdk-15.04 -s vivid maint dpkg -i --force-all *.deb
popd
