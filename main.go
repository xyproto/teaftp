package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pin/tftp"
	"github.com/sirupsen/logrus"
	"github.com/xyproto/env"
)

const versionString = "TeaFTP 2.0.0"

var (
	// allowed filename string prefixes or suffixes. Non-empty slices are put to use.
	allowedPrefixes []string
	allowedSuffixes []string
)

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

	allowed := true

	// If the allowedPrefixes slice is not empty, check if the filename has a valid suffix
	if len(allowedPrefixes) > 0 {
		allowed = false
		for _, prefix := range allowedPrefixes {
			if strings.HasPrefix(filename, prefix) {
				allowed = true
				break
			}
		}

	}

	// If the allowedSuffixes slice is not empty, check if the filename has a valid suffix
	if len(allowedSuffixes) > 0 {
		allowed = false
		for _, suffix := range allowedSuffixes {
			if strings.HasSuffix(filename, suffix) {
				allowed = true
				break
			}
		}
	}

	// if the allowedPrefixes slice is not empty, check if the filename has a valid prefix

	// Check if the read request is allowed, or not
	if !allowed {
		if remoteAddr != "" && localAddr != "" {
			logrus.Errorf("[DENIED] Read request of %s from %s to %s: prefix or suffix not allowed", filename, remoteAddr, localAddr)
		}
		return fmt.Errorf("%s does not have an allowed prefix or suffix", filename)
	}

	// Log the request
	if remoteAddr != "" && localAddr != "" {
		logrus.Infof("Read request of %s from %s to %s", filename, remoteAddr, localAddr)
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
			logrus.Infof("Write request of %s from %s to %s", filename, remoteAddr, localAddr)
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
	// Is the server read-only?
	readOnly := true

	// TODO: Use a proper command line argument handler

	//  Handle allowed filename suffixes
	if len(os.Args) > 1 {
		if os.Args[1] == "--help" {
			fmt.Println(versionString + `

Any given arguments that are not flags are interpreted as filename suffixes
that are added to the list of allowed filename suffixe.
Example, for serving only filenames ending with .txt:

teaftp ".txt"
`)
			os.Exit(0)
		} else if os.Args[1] == "--write" { // Undocumented feature
			logrus.Infoln("Enabled write mode")
			readOnly = false
		}
		allowedSuffixes = os.Args[1:]
		//allowedPrefixes = os.Args[1:]
	}

	fmt.Println(versionString + "\nSimple, read-only TFTP server")

	// use nil in place of handler to disable read or write operations
	s := tftp.NewServer(readHandler, genWriteHandler(readOnly))
	s.SetTimeout(5 * time.Second)                        // optional
	err := s.ListenAndServe(":" + env.Str("PORT", "69")) // blocks until s.Shutdown() is called
	if err != nil {
		logrus.Errorf("server: %s", err)
		os.Exit(1)
	}
}
