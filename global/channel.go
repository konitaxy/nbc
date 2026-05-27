package global

var (
	LogChannel = make(chan interface{}, 100)
)

func Push(log interface{}) {
	LogChannel <- log
}
