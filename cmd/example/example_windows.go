//go:build windows

package main

import (
	"github.com/Microsoft/go-winio"
)

func RunWithPrivileges(cb func() error) error {
	return winio.RunWithPrivilege("SeSystemEnvironmentPrivilege", cb)
}
