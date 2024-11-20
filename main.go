package main

func main() {
}

type MyErr struct {
	Msg string
}

func (e MyErr) Error() string {
	return e.Msg
}
