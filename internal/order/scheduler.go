package order

import "github.com/go-co-op/gocron"

type Scheduler interface {
	MapCron(*gocron.Scheduler)
}
