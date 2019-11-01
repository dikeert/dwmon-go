package sinks

func Stdout(output string) error {
	println(output)
	return nil
}
