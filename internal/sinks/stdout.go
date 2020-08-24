package sinks

import "os"

func Stdout(output string) error {
	os.Stdout.WriteString(output)
	return nil
}
