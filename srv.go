package main

import (
	"github.com/as/srv/fs"
)

var (
	srv               *fs.Server
	client            *fs.Client
	srverr, clienterr error
)

func createnetworks() (fatal error) {
	if *srvaddr != "" {
		srv, srverr = fs.Serve("tcp", *srvaddr)
	}
	if *dialaddr != "" {
		client, clienterr = fs.Dial("tcp", *dialaddr)
		if clienterr != nil {
			return clienterr
		}
	}
	if srv != nil {
		logf("listening for remote connections")
	}
	if client != nil {
		logf("connecting to remote filesystem")
	}
	return nil
}

func newfsclient() fs.Fs {
	if client == nil {
		return &fs.Local{}
	}
	if clienterr != nil {
		client, clienterr = fs.Dial("tcp", *dialaddr)
	}
	if clienterr == nil {
		return client
	}
	panic(clienterr)
}
