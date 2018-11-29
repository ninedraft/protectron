package heuristic

import (
	"time"
)

type Message struct {
	From       User
	IsRepost   bool
	RepostFrom RepostFrom
	Text       string
}

type User struct {
	JoinDate time.Time
	Username string
	Name     string
	ID       int64
}

type RepostFrom struct {
	ChatID       int64
	ChatUsername string
	ChatTitle    string
}

type Heuristic interface {
	Name() string
	Check(message Message) (p float64, err error)
}

func RoundP(p float64) float64 {
	if p < 0 {
		return 0
	}
	if p > 1 {
		return 1
	}
	return p
}

type Heuristics []Heuristic

func (heuristics Heuristics) Eval(message Message) (Vec, error) {
	var vec = make(Vec, 0, len(heuristics))
	for _, heuristic := range heuristics {
		var p, err = heuristic.Check(message)
		if err != nil {
			return nil, err
		}
		vec = append(vec, RoundP(p))
	}
	return vec, nil
}
