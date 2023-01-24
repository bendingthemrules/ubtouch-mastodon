# Mastodon

[![OpenStore](https://open-store.io/badges/en_US.png)](https://open-store.io/app/nl.btr.mastodon)

Mastodon with notifications

## Building

Build the application by executing:

```bash
./build.sh
```

Push to device with (make sure the versions match):

```bash
adb push nl.btr.mastodon_1.0.0_arm64.click
```

Open device with adb shell and execute:

```bash
pkcon install-local --allow-untrusted nl.btr.mastodon_1.0.0_arm64.click
```

### Push helper

Build pushHelper by executing:

```bash
cd pushnotifications/executable/ && qtdeploy build
```

## License

This project is licened under the GNU GPL v3 license
