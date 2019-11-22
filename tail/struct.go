package tail

// Follower describes an object that continuously emits a stream of lines
type Follower interface {
	Lines() chan string
	OnError(func(error))
}
