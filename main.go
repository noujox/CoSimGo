package main

import (
	"fmt"

	"github.com/agoussia/godes"
)

const (
	SHUT_DOWN_TIME = 1 * 60
	// ACTIVIDADES del truck
	PARADO        = 101
	CARGANDO      = 102
	TRANSPORTANDO = 103
	DESCARGANDO   = 104
	REGRESANDO    = 105
	AVERIADO      = 110
	// ESTADOS del truck
	FALLA     = 150
	OPERATIVO = 151
	//tipos de pile
	STOCK_PILE_TYPE = 201
	CHARGER_TYPE    = 202
)

var charger_time *godes.UniformDistr = godes.NewUniformDistr(true)
var truck_time *godes.UniformDistr = godes.NewUniformDistr(true)
var tim_gen *godes.UniformDistr = godes.NewUniformDistr(true)
var breaksGen *godes.ExpDistr = godes.NewExpDistr(true)

var pils dispatcher
var chars dispatcher

type truckMachine struct {
	*godes.Runner
	id, value int
	busy      *godes.BooleanControl
	activitie rune //Parado, Transito,Cargando, Descargando, Moviendo
}
type truck struct {
	*godes.Runner
	machine *truckMachine
}

func (trm truck) Run() {
	machine := trm.machine
	for {
		godes.Advance(breaksGen.Get(1 / 20))
		if machine.activitie == PARADO {
			break
		}

		interrupted := godes.GetSystemTime()
		godes.Interrupt(machine)
		fmt.Println("Truck ", machine.id, " fallo \t", godes.GetSystemTime())
		godes.Advance(5)
		if machine.activitie == PARADO {
			break
		}
		godes.Resume(machine, godes.GetSystemTime()-interrupted)

	}
}

func (tr truckMachine) Run() {
	for {
		if SHUT_DOWN_TIME < godes.GetSystemTime() {
			tr.activitie = PARADO
			break
		}

		switch tr.activitie {
		case PARADO:
			fmt.Println("T ", tr.id, ": Parado..  \t", godes.GetSystemTime())
			godes.Advance(tim_gen.Get(1, 5))
			tr.activitie = CARGANDO

		case CARGANDO:
			fmt.Println("T ", tr.id, ": Cargando \t", godes.GetSystemTime())
			chars.lis.Front().Value.(shovel).qe.Place(tr)
			chars.lis.Front().Value.(shovel).empty_qe.Set(false)
			chars.nextList()
			tr.busy.Set(true)
			tr.busy.Wait(false)
			tr.activitie = TRANSPORTANDO

		case TRANSPORTANDO:
			fmt.Println("T ", tr.id, ": Transportando\t", godes.GetSystemTime())
			godes.Advance(tim_gen.Get(5, 10))
			tr.activitie = DESCARGANDO

		case DESCARGANDO:
			fmt.Println("T ", tr.id, ": Descargando\t", godes.GetSystemTime())
			pils.lis.Front().Value.(stockPile).qe.Place(tr)
			pils.lis.Front().Value.(stockPile).empty_qe.Set(false)
			pils.nextList()
			tr.busy.Set(true)
			tr.busy.Wait(false)
			tr.activitie = REGRESANDO

		case REGRESANDO:
			fmt.Println("T ", tr.id, ": Volviendo \t", godes.GetSystemTime())
			godes.Advance(tim_gen.Get(1, 5))
			tr.activitie = CARGANDO

		default:
			fmt.Println("explotoo")
		}
	}
}

func (tr truckMachine) receive(x int) bool {
	if tr.activitie == CARGANDO {
		tr.value = x
		return true
	} else {
		fmt.Println("Truck can't receive the payload")
		return false
	}
}

func (tr truckMachine) get() int {
	return tr.value
}

func newTruck(id int) *truck {
	tm := &truckMachine{&godes.Runner{}, id, 0, godes.NewBooleanControl(), PARADO}
	godes.AddRunner(tm)
	return &truck{&godes.Runner{}, tm}
}

func main() {
	pils = pils.init(STOCK_PILE_TYPE)
	chars = chars.init(CHARGER_TYPE)

	godes.AddRunner(newTruck(1))
	godes.AddRunner(newTruck(2))

	for i := 0; i < 3; i++ {
		pils.addList(stockPile{i, 3, godes.NewFIFOQueue("pile"), godes.NewBooleanControl(), []pilePort{}})
		pils.lis.Back().Value.(stockPile).init()
	}
	for i := 0; i < 3; i++ {
		chars.addList(shovel{&godes.Runner{}, i, godes.NewFIFOQueue("chars"), godes.NewBooleanControl()})
		godes.AddRunner(chars.lis.Back().Value.(shovel))
	}
	godes.Run()
	godes.WaitUntilDone()
}
