package main
import "github.com/agoussia/godes"

type charger struct {
	*godes.Runner
	id       int
	qe       *godes.FIFOQueue
	empty_qe *godes.BooleanControl
}

func (ch charger) Run() {
	var x int
	ch.empty_qe.Set(true)
	for {
		if SHUT_DOWN_TIME < godes.GetSystemTime() {
			break
		}
		ch.empty_qe.Wait(false)
		tru := ch.qe.Get().(truck)
		if ch.qe.Len() == 0 {
			ch.empty_qe.Set(true)
		}
		godes.Yield()

		tru.receive(x)
		tru.busy.Set(false)
		x++

	}
}