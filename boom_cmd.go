package megaboom

import (
	"strconv"
)

type boomCommand struct {
	total       int
	concurrency int
	endpoint    string
}

func (b boomCommand) Slice() []string {
	return []string{"boom", "-c", strconv.Itoa(b.concurrency), "-n", strconv.Itoa(b.total), b.endpoint}

}
