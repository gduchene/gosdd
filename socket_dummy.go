// +build !linux nosystemd

package gosdd

import "os"

func sdListenFDs(bool) ([]*os.File, error) {
	return nil, nil
}

func sdListenFDsWithNames(bool) (map[string]*os.File, error) {
	return nil, nil
}
