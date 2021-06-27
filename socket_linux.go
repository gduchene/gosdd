// SPDX-FileCopyrightText: © 2021 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

// +build linux,!nosystemd

package gosdd

// #cgo LDFLAGS: -lsystemd
// #include <stdlib.h>
// #include <string.h>
// #include <systemd/sd-daemon.h>
import "C"

import (
	"fmt"
	"os"
	"unsafe"
)

func sdListenFDs(unsetenv bool) ([]*os.File, error) {
	i := C.int(0)
	if unsetenv {
		i = C.int(1)
	}
	c := C.sd_listen_fds(i)
	if c < 0 {
		return nil, fmt.Errorf("sd_listen_fds: %s", C.GoString(C.strerror(-c)))
	}
	if c == 0 {
		return nil, nil
	}
	fds := make([]*os.File, 0, c)
	for fd := uintptr(C.SD_LISTEN_FDS_START); fd < uintptr(C.SD_LISTEN_FDS_START+c); fd++ {
		fds = append(fds, os.NewFile(fd, ""))
	}
	return fds, nil
}

func sdListenFDsWithNames(unsetenv bool) (map[string]*os.File, error) {
	i := C.int(0)
	if unsetenv {
		i = C.int(1)
	}
	var arr **C.char
	c := C.sd_listen_fds_with_names(i, &arr)
	if c < 0 {
		return nil, fmt.Errorf("sd_listen_fds_with_names: %s", C.GoString(C.strerror(-c)))
	}
	if c == 0 {
		return nil, nil
	}
	// See https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices.
	names := (*[1 << 28]*C.char)(unsafe.Pointer(arr))[:c:c]
	fds := make(map[string]*os.File)
	for fd := uintptr(C.SD_LISTEN_FDS_START); fd < uintptr(C.SD_LISTEN_FDS_START+c); fd++ {
		name := C.GoString(names[int(fd-C.SD_LISTEN_FDS_START)])
		fds[name] = os.NewFile(fd, name)
		C.free(unsafe.Pointer(names[int(fd-C.SD_LISTEN_FDS_START)]))
	}
	C.free(unsafe.Pointer(arr))
	return fds, nil
}
