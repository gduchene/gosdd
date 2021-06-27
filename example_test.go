// SPDX-FileCopyrightText: © 2021 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

package gosdd_test

import (
	"log"
	"net"
	"net/http"

	"go.awhk.org/gosdd"
)

// This example gets the file descriptors passed to the process by
// systemd and starts a server using it.
func ExampleSDListenFDs() {
	fds, err := gosdd.SDListenFDs(true)
	if err != nil {
		log.Fatalf("Error while getting file descriptors: %s.", err)
	}
	if len(fds) != 1 {
		log.Fatalln("Exactly one file descriptor can be handled.")
	}

	ln, err := net.FileListener(fds[0])
	if err != nil {
		log.Fatalf("Failed to create listener: %s.", err)
	}
	srv := http.Server{Handler: http.FileServer(http.Dir("/tmp"))}
	log.Println(srv.Serve(ln))
}
