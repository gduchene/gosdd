// SPDX-FileCopyrightText: © 2021 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

// This implements a simple example that can be tested on a machine
// running systemd.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"go.awhk.org/gosdd"
)

var useNames = flag.Bool("use-names", false, "whether to use SDListenFDsWithNames or not")

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)

	if !*useNames {
		fds, err := gosdd.SDListenFDs(true)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Received %d file descriptors from systemd.", len(fds))
		listenAll(fds)
		return
	}

	namedFDs, err := gosdd.SDListenFDsWithNames(true)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Received %d named file descriptors from systemd.", len(namedFDs))
	fds := make([]*os.File, 0, len(namedFDs))
	for name, fd := range namedFDs {
		log.Printf("Adding %q.", name)
		fds = append(fds, fd)
	}
	listenAll(fds)
}

func listenAll(fds []*os.File) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	srvs := make([]echoServer, 0, len(fds))
	wg := &sync.WaitGroup{}
	for _, fd := range fds {
		ln, err := net.FileListener(fd)
		if err != nil {
			log.Printf("Failed to make a listener: %s.", err)
			continue
		}
		srv := echoServer{ln, wg}
		srvs = append(srvs, srv)
		go srv.start()
	}

	<-sig
	for _, srv := range srvs {
		srv.stop()
	}
	wg.Wait()
}

type echoServer struct {
	ln net.Listener
	wg *sync.WaitGroup
}

func (*echoServer) handle(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println(err)
		}
	}()

	log.Printf("New connection, local address is %s, remote address is %s.", conn.LocalAddr(), conn.RemoteAddr())
	if _, err := fmt.Fprintln(conn, "Hello World!"); err != nil {
		log.Println(err)
		return
	}
	r := bufio.NewReader(conn)
	for {
		s, err := r.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		s = strings.Trim(s, "\n ")
		log.Printf("Received %q.", s)
		if _, err := fmt.Fprintf(conn, "You said %q!\n", s); err != nil {
			log.Println(err)
			return
		}
	}
}

func (srv *echoServer) start() {
	defer srv.wg.Done()
	srv.wg.Add(1)

	for {
		conn, err := srv.ln.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		go srv.handle(conn)
	}
}

func (srv *echoServer) stop() {
	if err := srv.ln.Close(); err != nil {
		log.Println(err)
	}
}
