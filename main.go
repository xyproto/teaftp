package main

import (
	"fmt"
	"github.com/pin/tftp"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

const versionString = "TeaFTP 1.0"

// readHandler is called when client starts file download from server
func readHandler(filename string, rf io.ReaderFrom) error {
	remoteAddr := ""
	if raddr, ok := rf.(tftp.OutgoingTransfer); ok {
		r := raddr.RemoteAddr()
		remoteAddr = r.String()
	}
	localAddr := ""
	if laddr, ok := rf.(tftp.RequestPacketInfo); ok {
		localAddr = laddr.LocalIP().String()
	}
	if remoteAddr != "" && localAddr != "" {
		logrus.Infof("Read request from %s to %s", remoteAddr, localAddr)
	}

	file, err := os.Open(filename)
	if err != nil {
		logrus.Error(err)
		return err
	}

	// Find the size of the file
	fi, err := file.Stat()
	if err != nil {
		// Could not obtain stat, handle error
		logrus.Error(err)
		return err
	}
	fileSize := fi.Size()

	// Set transfer size before calling ReadFrom.
	rf.(tftp.OutgoingTransfer).SetSize(fileSize)

	n, err := rf.ReadFrom(file)
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("%d bytes sent", n)
	return nil
}

func genWriteHandler(readOnly bool) func(string, io.WriterTo) error {
	// writeHandler is called when client starts file upload to server
	return func(filename string, wt io.WriterTo) error {
		remoteAddr := ""
		if raddr, ok := wt.(tftp.IncomingTransfer); ok {
			r := raddr.RemoteAddr()
			remoteAddr = r.String()
		}
		localAddr := ""
		if laddr, ok := wt.(tftp.RequestPacketInfo); ok {
			localAddr = laddr.LocalIP().String()
		}
		if remoteAddr != "" && localAddr != "" {
			logrus.Infof("Write request from %s to %s", remoteAddr, localAddr)
		}

		if readOnly {
			logrus.Warnf("Client wants to upload %s, but server is read-only", filename)
			filename = os.DevNull
		}
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0664)
		if err != nil {
			logrus.Error(err)
			return err
		}
		n, err := wt.WriteTo(file)
		if err != nil {
			// An error is expected when trying to write to /dev/null,
			// since the file already exists.
			if !readOnly {
				logrus.Error(err)
			}
			return err
		}
		logrus.Infof("%d bytes received", n)
		return nil
	}
}

func main() {
	fmt.Println(versionString + "\nSimple, read-only TFTP server")
	// Is the server read-only?
	readOnly := true
	// use nil in place of handler to disable read or write operations
	s := tftp.NewServer(readHandler, genWriteHandler(readOnly))
	s.SetTimeout(5 * time.Second)  // optional
	err := s.ListenAndServe(":69") // blocks until s.Shutdown() is called
	if err != nil {
		logrus.Errorf("server: %s", err)
		os.Exit(1)
	}
}