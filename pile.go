package main

import "github.com/agoussia/godes"

type stockPile struct {
	id, nPorts int
	qe         *godes.FIFOQueue
	empty_qe   *godes.BooleanControl
	ports      []pilePort
}

func (sp stockPile) init() stockPile {
	for i := 0; i < sp.nPorts; i++ {
		sp.ports = append(sp.ports, (pilePort{&godes.Runner{}, i, &sp.empty_qe, &sp.qe}))
		godes.AddRunner(sp.ports[i])
	}
	return sp
}

type pilePort struct {
	*godes.Runner
	id       int
	empty_qe **godes.BooleanControl
	qe       **godes.FIFOQueue
}

func (pp pilePort) Run() {
	pp.empty_qe.Set(true)
	for {
		if SHUT_DOWN_TIME < godes.GetSystemTime() {
			break
		}
		pp.empty_qe.Wait(false)
		tru := pp.qe.Get().(truckMachine)
		if pp.qe.Len() == 0 {
			pp.empty_qe.Set(true)
		}
		godes.Yield()

		tru.get()
		tru.busy.Set(false)

	}
}
