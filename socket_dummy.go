// SPDX-FileCopyrightText: © 2021 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

//go:build !linux || nosystemd

package gosdd

import "os"

func sdListenFDs(bool) ([]*os.File, error) {
	return nil, ErrNoSDSupport
}

func sdListenFDsWithNames(bool) (map[string]*os.File, error) {
	return nil, ErrNoSDSupport
}
