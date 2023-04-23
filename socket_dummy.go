// SPDX-FileCopyrightText: © 2021 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

//go:build !linux || nosystemd

package gosdd

import "os"

func SDListenFDs(bool) ([]*os.File, error) {
	return nil, ErrNoSDSupport
}

func SDListenFDsWithNames(bool) (map[string]*os.File, error) {
	return nil, ErrNoSDSupport
}
