package main

import "github.com/agoussia/godes"

type pile struct {
	*godes.Runner
	qe       *godes.FIFOQueue
	empty_qe *godes.BooleanControl
}

func (pl pile) Run() {
	pl.empty_qe.Set(true)
	for {
		if SHUT_DOWN_TIME < godes.GetSystemTime() {
			break
		}
		pl.empty_qe.Wait(false)
		tru := pl.qe.Get().(truck)
		if pl.qe.Len() == 0 {
			pl.empty_qe.Set(true)
		}
		godes.Yield()

		tru.get()
		tru.busy.Set(false)

	}
}