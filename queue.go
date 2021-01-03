package main

type Queue struct {
	Tracklist    []string `json:"tracklist"`
	CurrentIndex int      `json:"current_index"`
}

func NewQueue() *Queue {
	return &Queue{
		Tracklist:    make([]string, 0, 64),
		CurrentIndex: 0,
	}
}

func (q *Queue) JumpToIndex(index int) {
	q.CurrentIndex = index
}

func (q *Queue) AddToQueue(trackName string) {
	q.Tracklist = append(q.Tracklist, trackName)
}

func (q *Queue) NextTrack() {
	q.CurrentIndex += 1
}

func (q *Queue) PrevTrack() {
	q.CurrentIndex -= 1
}
