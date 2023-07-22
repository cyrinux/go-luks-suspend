// Copyright (c) 2017 Sung Pae <self@sungpae.com>
// Distributed under the MIT license.
// http://www.opensource.org/licenses/mit-license.php

// +build !windows

package sys

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// SetTTYIoctl identifies the three tcsetattr ioctls described in ioctl_tty(2)
type SetTTYIoctl uintptr

// These are typed for compile time safety.
const (
	TCSETS  SetTTYIoctl = unix.TCSETS
	TCSETSW             = unix.TCSETSW
	TCSETSF             = unix.TCSETSF
)

// GetTTYState writes the TTY state of fd to termios.
func GetTTYState(fd uintptr, termios *syscall.Termios) error {
	_, _, err := Ioctl(fd, syscall.TCGETS, uintptr(unsafe.Pointer(termios)))
	return err
}

// SetTTYState alters the TTY state of fd to match termios.
func SetTTYState(fd uintptr, action SetTTYIoctl, termios *syscall.Termios) error {
	// tcsetattr(3):
	// Note that tcsetattr() returns success if any of the requested changes
	// could be successfully carried out. Therefore, when making multiple changes
	// it may be necessary to follow this call with a further call to tcgetattr()
	// to check that all changes have been performed successfully.
	state := syscall.Termios{}
	for {
		_, _, err := Ioctl(fd, uintptr(action), uintptr(unsafe.Pointer(termios)))
		if err != nil {
			return err
		}

		if err := GetTTYState(fd, &state); err != nil {
			return err
		}

		if state == *termios {
			return nil
		}
	}
}

// AlterTTY provides a simple way to change a TTY and restore it later.
//
// Function f receives a copy of the current TTY state of fd and modifies it.
// The TTY indicated by fd is then changed to match this modified state.
//
// A function is returned that will return the TTY to its original state if it
// was altered. If the TTY was not altered, restoreTTY will be nil.
func AlterTTY(fd uintptr, action SetTTYIoctl, f func(*syscall.Termios)) (restoreTTY func() error, err error) {
	oldstate := syscall.Termios{}

	if err := GetTTYState(fd, &oldstate); err != nil {
		return nil, err
	}

	restoreTTY = func() error { return SetTTYState(fd, action, &oldstate) }

	newstate := oldstate
	f(&newstate)

	return restoreTTY, SetTTYState(fd, action, &newstate)
}
