//go:build !darwin && !linux

package main

func fdIsTTY(uintptr) bool { return false }
