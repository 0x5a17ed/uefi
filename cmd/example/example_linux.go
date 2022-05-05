package main

func RunWithPrivileges(cb func() error) error {
	return cb()
}
