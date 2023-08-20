<img src="img/teaftp.svg" width="128">

# TeaFTP

![Build](https://github.com/xyproto/teaftp/workflows/Build/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/teaftp)](https://goreportcard.com/report/github.com/xyproto/teaftp) [![License](https://img.shields.io/badge/license-BSD-green.svg?style=flat)](https://raw.githubusercontent.com/xyproto/teaftp/main/LICENSE)

Simple, read-only TFTP server.

### Features and Limitations

* Suitable for dealing with hardware devices that read files over TFTP at boot (PXE).
* Security is provided by providing a list of whitelisted prefixes, suffixes and/or running the server from within a Docker container.
* If not whitelisted filename prefixes or suffixes are provided, the server may share ANY file on the system, but not write to anything.
  * Consider using the provided Docker container as a method to serve only a select group of files.
  * Alternatively, utilize the list of allowed prefixes or suffixes for added security.
* Every access is logged to stdout.

### Requirements

    Go 1.17 or later

### Installation with Go 1.17 or Later

    go install github.com/xyproto/teaftp@latest

### Running

#### Directly

Navigate to the directory where you intend to share files:

With sudo:

    sudo ./teaftp

On Linux, you can `teaftp` it to `/usr/bin` and grant additional capabilities using `setcap`:

    sudo install -Dm755 teaftp /usr/bin/teaftp
    sudo setcap cap_net_bind_service=+ep /usr/bin/teaftp

Starting the server:

    teaftp

#### Docker

To build the Docker container and copy the contents of the `static` directory to `/srv/tftp` inside the container:

    docker build . -t teaftp

To run TeaFTP with Docker:

    docker run --network=host -t teaftp

To run TeaFTP with Docker and serve on port 9000 instead of port 69:

    docker run -ePORT=9000 --network=host -t teaftp

#### Allowed Suffixes

You can pass allowed filename suffixes as arguments to TeaFTP. When no arguments are given, there's no restriction on the file suffixes.

Example:

    sudo ./teaftp ".txt"

This configuration will only serve filenames that end with `.txt`.

### Dependencies

* [pin/tftp](https://github.com/pin/tftp)
* [sirupsen/logrus](https://github.com/sirupsen/logrus)
* [urfave/cli](https://github.com/urfave/cli)

### General information

* Version: 1.3.0
* License: BSD-3
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
