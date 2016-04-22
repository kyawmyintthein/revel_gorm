package main

import "fmt"

type flagValue string


func (d *flagValue) String() string {
	return fmt.Sprint(*d)
}

func (d *flagValue) Set(value string) error {
	*d = flagValue(value)
	return nil
}

