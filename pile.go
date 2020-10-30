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
		sp.ports = append(sp.ports, (pilePort{&godes.Runner{}, i, &sp}))
		godes.AddRunner(sp.ports[i])
	}
	return sp
}

type pilePort struct {
	*godes.Runner
	id  int
	spP *stockPile
}

func (pp pilePort) Run() {
	pp.spP.empty_qe.Set(true)
	for {
		if SHUT_DOWN_TIME < godes.GetSystemTime() {
			break
		}
		pp.spP.empty_qe.Wait(false)
		tru := pp.spP.qe.Get().(truckMachine)
		if pp.spP.qe.Len() == 0 {
			pp.spP.empty_qe.Set(true)
		}
		godes.Yield()

		tru.get()
		tru.busy.Set(false)

	}
}
