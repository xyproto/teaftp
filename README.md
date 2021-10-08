# TeaFTP [![Build Status](https://travis-ci.com/xyproto/teaftp.svg?branch=master)](https://travis-ci.com/xyproto/teaftp)

Simple, read-only TFTP server.

* Will happily share ANY file on the system, but does not have access to write to any file.
  * Use the provided Docker container for a way to serve only a limited selection of files.
  * Or use the list of allowed prefixes or suffixes, as described below.
* TeaFTP may be suitable for dealing with hardware devices that read files over TFTP at boot.
* Every access is logged to stdout.

### Requirements

    Go >= 1.11

### Installation with Go >= 1.17

    go install github.com/xyproto/teaftp@latest

### Running

#### Directly

In the directory where you wish to share files:

Either:

    sudo ./teaftp

Or as root or with the correct Linux capabilities:

    ./teaftp

Here is one way of letting TeaFTP listen to the TFTP port without running as root, on Linux:

    sudo setcap cap_net_bind_service=+ep /usr/bin/teaftp

#### Docker

To build the Docker container, and also copy in the contents of the "static" directory to /srv/tftp within the container:

    docker build . -t teaftp

To run TeaFTP with Docker:

    docker run --network=host -ti teaftp

To run TeaFTP with Docker and serve on port 9000:

    docker run -ePORT=9000 --network=host -ti teaftp

#### Allowed suffixes

Any arguments given to TeaFTP are added to the list of allowed filename suffixes. If no arguments are given, the list of allowed suffixes is not in use.

Example:

    sudo ./teaftp ".txt"

This only serves filenames ending with `.txt`.

### Uses

* [pin/tftp](https://github.com/pin/tftp)
* [sirupsen/logrus](https://github.com/sirupsen/logrus)

### General info

* Version: 2.0.0
* License: BSD-3
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
