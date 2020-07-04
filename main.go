package main

import (
	"container/list"
	"fmt"

	"github.com/agoussia/godes"
)

const (
	SHUT_DOWN_TIME = 1 * 60
)

var charger_time *godes.UniformDistr = godes.NewUniformDistr(true)
var truck_time *godes.UniformDistr = godes.NewUniformDistr(true)
var tim_gen *godes.UniformDistr = godes.NewUniformDistr(true)

var pils dispatcher
var chars dispatcher

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

type pile struct {
	*godes.Runner
	id       int
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

type truck struct {
	*godes.Runner
	id, value int
	busy      *godes.BooleanControl
	state     rune //Parado, Transito,Cargando, Descargando, Moviendo
}

func (tr truck) Run() {
	for {
		if SHUT_DOWN_TIME < godes.GetSystemTime() {
			tr.state = 'P'
			break
		}

		switch tr.state {
		case 'P':
			fmt.Println("T ", tr.id, " : Parado... iniciando")
			godes.Advance(tim_gen.Get(1, 5))
			tr.state = 'C'
		case 'C':
			fmt.Println("T ", tr.id, " : Cargando")

			chars.lis.Front().Value.(charger).qe.Place(tr)
			chars.lis.Front().Value.(charger).empty_qe.Set(false)
			chars.nextList()
			tr.busy.Set(true)
			tr.busy.Wait(false)
			tr.state = 'T'
		case 'T':
			fmt.Println("T ", tr.id, " : Transportando")
			godes.Advance(tim_gen.Get(5, 10))
			tr.state = 'D'
		case 'D':
			fmt.Println("T ", tr.id, " : Descargando")
			pils.lis.Front().Value.(pile).qe.Place(tr)
			pils.lis.Front().Value.(pile).empty_qe.Set(false)
			pils.nextList()
			tr.busy.Set(true)
			tr.busy.Wait(false)
			tr.state = 'M'
		case 'M':
			fmt.Println("T ", tr.id, " : Volviendo a cargar")
			godes.Advance(tim_gen.Get(1, 5))
			tr.state = 'C'
		default:
			fmt.Println("exploto")
		}
	}
}
func (tr truck) receive(x int) bool {
	if tr.state == 'C' {
		tr.value = x
		return true
	} else {
		fmt.Println("Truck can't receive the payload")
		return false
	}
}
func (tr truck) get() int {
	return tr.value
}

type dispatcher struct {
	tipe rune
	lis  *list.List
}

func (ds dispatcher) init(r rune) dispatcher {
	return dispatcher{r, list.New()}
}

func (ds dispatcher) addList(x interface{}) {
	ds.lis.PushBack(x)
}
func (ds dispatcher) nextList() *list.Element {
	e := ds.lis.Front()
	ds.lis.MoveToBack(e)
	return e
}

/* func (ds dispatcher) dispatch() interface{} {
	//error
	return ds.e
} */

func main() {
	pils = pils.init('p')
	chars = chars.init('c')
	godes.AddRunner(&truck{&godes.Runner{}, 1, 0, godes.NewBooleanControl(), 'P'})
	godes.AddRunner(&truck{&godes.Runner{}, 2, 0, godes.NewBooleanControl(), 'P'})
	for i := 0; i < 3; i++ {
		pils.addList(pile{&godes.Runner{}, i, godes.NewFIFOQueue("pile"), godes.NewBooleanControl()})
	}
	for i := 0; i < 3; i++ {
		chars.addList(charger{&godes.Runner{}, i, godes.NewFIFOQueue("chars"), godes.NewBooleanControl()})
	}
	godes.Run()
	godes.WaitUntilDone()
}
