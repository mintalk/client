package network

import "time"

func ScheduleRepeatedTask(interval time.Duration, event func()) {
	go func() {
		for {
			time.Sleep(interval)
			event()
		}
	}()
}
