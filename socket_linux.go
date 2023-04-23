// SPDX-FileCopyrightText: © 2021 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

//go:build linux && !nosystemd

package gosdd

// #cgo LDFLAGS: -lsystemd
// #include <stdlib.h>
// #include <string.h>
// #include <systemd/sd-daemon.h>
import "C"

import (
	"os"
	"syscall"
	"unsafe"
)

func SDListenFDs(unsetenv bool) ([]*os.File, error) {
	i := C.int(0)
	if unsetenv {
		i = C.int(1)
	}
	c := C.sd_listen_fds(i)
	if c < 0 {
		return nil, syscall.Errno(-c)
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

func SDListenFDsWithNames(unsetenv bool) (map[string]*os.File, error) {
	i := C.int(0)
	if unsetenv {
		i = C.int(1)
	}
	var arr **C.char
	c := C.sd_listen_fds_with_names(i, &arr)
	if c < 0 {
		return nil, syscall.Errno(-c)
	}
	if c == 0 {
		return nil, nil
	}
	names := unsafe.Slice(arr, c)
	fds := make(map[string]*os.File)
	for fd := uintptr(C.SD_LISTEN_FDS_START); fd < uintptr(C.SD_LISTEN_FDS_START+c); fd++ {
		name := C.GoString(names[int(fd-C.SD_LISTEN_FDS_START)])
		fds[name] = os.NewFile(fd, name)
		C.free(unsafe.Pointer(names[int(fd-C.SD_LISTEN_FDS_START)]))
	}
	C.free(unsafe.Pointer(arr))
	return fds, nil
}
