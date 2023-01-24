#!/bin/bash
set -e

# Build package
qtdeploy build
mkdir -p build

cp manifest.json build
cp mastodon.apparmor build
cp mastodon.desktop build
cp mastodon.url-dispatcher.json build

cp deploy/linux/mastodon-client build
mv build/mastodon-client build/mastodon

# Push notifications helper
cd pushnotifications/executable/ && qtdeploy build
cd ../..
cp pushnotifications/executable/deploy/linux/executable build
mv build/executable build/mastodonHelper
cp pushnotifications/mastodonHelper.apparmor.json build
cp pushnotifications/mastodonHelper.json build

cp -R assets build

cd build && click build .
