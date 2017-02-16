package notify

var S Sender

func init() {
	S = &Sender{make(chan string, 1024)}
}

type Sender struct {
	StringChan chan string
}
