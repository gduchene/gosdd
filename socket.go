// SPDX-FileCopyrightText: © 2021 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

// Package gosdd provides simple wrappers around useful functions
// provided by systemd.
//
// On systems that are not Linux, or if the nosystemd build tag is set,
// safe defaults are returned: zero or nil values, and no error will be
// returned.
//
// Reference
//
// https://www.freedesktop.org/software/systemd/man/sd_listen_fds.html
// is the documentation for the C API.
package gosdd

import (
	"errors"
	"os"
)

// ErrNoSDSupport is a generic error that is returned when gosdd has no
// systemd support, either because the library is compiled on a system
// that is not Linux or because it was explicitly disabled with the
// ‘nosystemd’ build tag.
var ErrNoSDSupport = errors.New("no systemd support")

// SDListenFDs is a wrapper around sd_listen_fds.
func SDListenFDs(unsetenv bool) ([]*os.File, error) {
	return sdListenFDs(unsetenv)
}

// SDListenFDsWithNames is a wrapper around sd_listen_fds_with_names.
func SDListenFDsWithNames(unsetenv bool) (map[string]*os.File, error) {
	return sdListenFDsWithNames(unsetenv)
}
