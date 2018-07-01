package main

import "testing"

func TestStartChain(t *testing.T) {
	errs := StartChain(func() error {
		println("started")
		return nil
	}).Next(func() error {
		println("1 thread")
		return nil
	}).Parallel(func() error {
		println("2 thread")
		return nil
	}).Parallel(func() error {
		println("3 thread")
		return nil
	}).Parallel(func() error {
		println("4 thread")
		return nil
	}).Next(func() error {
		println("somewhere in the middle")
		return nil
	}).Next(func() error {
		println("1 thread after middle")
		return nil
	}).Parallel(func() error {
		println("2 thread after middle")
		return nil
	}).Next(func() error {
		println("the end!")
		return nil
	}).Run()
	if len(errs) != 0 {
		t.Fatal(errs)
	}
}