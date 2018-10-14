# TeaFTP [![Build Status](https://travis-ci.org/xyproto/teaftp.svg?branch=master)](https://travis-ci.org/xyproto/teaftp)

Simple, read-only TFTP server.

![teaftp](img/teaftp.gif)

* Will happily share ANY file on the system, but does not have acccess to write to any file.
* This may be suitable when dealing with hardware devices that read files over TFTP at boot.
* This is not suitable for running an online server.
* Every access is logged.

### Requirements

    Go >= 1.11

### Installation

    go get github.com/xyproto/teaftp

### Running

In the directory where you wish to share files:

Either:

    sudo ./teaftp

Or as root or with the correct Linux capabilities:

    ./teaftp

### License

MIT

### Uses

* [pin/tftp](https://github.com/pin/tftp)
* [sirupsen/logrus](https://github.com/sirupsen/logrus)

### Version

1.0.0
