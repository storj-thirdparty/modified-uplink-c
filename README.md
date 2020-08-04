# modified-uplink-c

Modified C library to be suggested for Storj V3 Network.

![Beta](https://img.shields.io/badge/version-beta-green.svg)

<img src="https://github.com/storj/storj/raw/master/resources/logo.png" width="100">

Storj is building a decentralized cloud storage network.
[Check out our white paper for more info!](https://storj.io/white-paper)

----

Storj is an S3-compatible platform and suite of decentralized applications that
allows you to store data in a secure and decentralized manner. Your files are
encrypted, broken into little pieces and stored in a global decentralized
network of computers. Luckily, we also support allowing you (and only you) to
retrieve those files!

# Build

Download and install the latest release of Go (at least Go 1.13) at [golang.org](https://golang.org/).

When ready, building the shared library is as easy as executing:

```
make build
```

The output is in the `.build` folder.

If you specifically need GPLv2 compatibility, you can use `make build-gpl2` instead,
which will compile the library without any Apache v2 dependencies (sadly, Apache v2
is incompatible with the GPLv2). Currently this results in slower hashing
performance (no github.com/minio/sha256-simd) and reduced debugging and
analysis infrastructure.

# License

This library is distributed under the
[MIT/expat](https://opensource.org/licenses/MIT) license.

# Support

If you have any questions or suggestions please reach out to us on
[our community forum](https://forum.storj.io/) or
email us at support@tardigrade.io.
