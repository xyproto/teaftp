package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pin/tftp"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const versionString = "TeaFTP 1.3.0"

var (
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
	app := &cli.App{
		Name:    "teaftp",
		Version: versionString,
		Usage:   "Simple, read-only TFTP server",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:   "write",
				Usage:  "Enable write mode (undocumented feature)",
				Hidden: true,
			},
			&cli.StringFlag{
				Name:    "port",
				Value:   "69",
				Usage:   "Port number for the TFTP server",
				EnvVars: []string{"PORT"}, // Backwards compatibility with original Docker example
			},
		},
		ArgsUsage: "[allowed suffixes]",
		Action: func(c *cli.Context) error {
			readOnly := !c.Bool("write")
			if c.Args().Present() {
				allowedSuffixes = c.Args().Slice()
			}

			// Your TFTP server setup and run code goes here
			s := tftp.NewServer(readHandler, genWriteHandler(readOnly))
			s.SetTimeout(5 * time.Second)
			addr := ":" + c.String("port")
			fmt.Println("Serving tea at localhost" + addr)
			err := s.ListenAndServe(addr)
			if err != nil {
				logrus.Errorf("server: %s", err)
				os.Exit(1)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
