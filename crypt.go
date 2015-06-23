// +build darwin freebsd netbsd

// Package crypt provides wrappers around functions available in crypt.h
//
// It wraps around the GNU specific extension (crypt) when the reentrant version
// (crypt_r) is unavailable. The non-reentrant version is guarded by a global lock
// so as to be safely callable from concurrent goroutines.
package crypt

import (
	"fmt"
	"sync"
)

/*
#define _XOPEN_SOURCE 700
#include <unistd.h>
*/
import "C"

var (
	mu sync.Mutex
)

// Crypt provides a wrapper around the glibc crypt() function.
// For the meaning of the arguments, refer to the README.
func Crypt(pass, salt string) (string, error) {
	c_pass := C.CString(pass)
	defer C.free(unsafe.Pointer(c_pass))

	c_salt := C.CString(salt)
	defer C.free(unsafe.Pointer(c_salt))

	mu.Lock()
	c_enc, err := C.crypt(c_pass, c_salt)
	defer C.free(unsafe.Pointer(c_enc))
	mu.Unlock()

	if c_enc == nil {
		return "", err
	}
	// Return nil error if the string is non-nil.
	// This happens because crypt seems to leak a spurious ENOENT which
	// is left over after it checks the /proc/sys file for fips mode.
	fmt.Println("Returning err %v", err)
	return C.GoString(c_enc), err
}
