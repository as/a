package main

import "github.com/as/srv/fs"

var (
	srv               *fs.Server
	client            *fs.Client
	srverr, clienterr error
)

func createnetworks() (fatal error) {
	if *srvaddr != "" {
		srv, srverr = fs.Serve("tcp", *srvaddr)
	}
	if *clientaddr != "" {
		client, clienterr = fs.Dial("tcp", *clientaddr)
		if clienterr != nil {
			return clienterr
		}
	}
	return nil
}

func newfsclient() fs.Fs {
	if client != nil {
		return client
	}
	return &fs.Local{}
}
