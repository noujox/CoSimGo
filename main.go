package main

import (
	"fmt"

	"github.com/agoussia/godes"
)

const (
	SHUT_DOWN_TIME = 1 * 60
	// estados del truck
	PARADO        = 101
	CARGANDO      = 102
	TRANSPORTANDO = 103
	DESCARGANDO   = 104
	REGRESANDO    = 105
	//tipos de pile
	STOCK_PILE_TYPE = 201
	CHARGER_TYPE    = 202
)

var charger_time *godes.UniformDistr = godes.NewUniformDistr(true)
var truck_time *godes.UniformDistr = godes.NewUniformDistr(true)
var tim_gen *godes.UniformDistr = godes.NewUniformDistr(true)

var pils dispatcher
var chars dispatcher

type truck struct {
	*godes.Runner
	id, value int
	busy      *godes.BooleanControl
	state     rune //Parado, Transito,Cargando, Descargando, Moviendo
}

func (tr truck) Run() {
	for {
		if SHUT_DOWN_TIME < godes.GetSystemTime() {
			tr.state = PARADO
			break
		}

		switch tr.state {
		case 'P':
			fmt.Println("T ", tr.id, " : Parado... iniciando")
			godes.Advance(tim_gen.Get(1, 5))
			tr.state = PARADO
		case 'C':
			fmt.Println("T ", tr.id, " : Cargando")

			chars.lis.Front().Value.(shovel).qe.Place(tr)
			chars.lis.Front().Value.(shovel).empty_qe.Set(false)
			chars.nextList()
			tr.busy.Set(true)
			tr.busy.Wait(false)
			tr.state = CARGANDO
		case 'T':
			fmt.Println("T ", tr.id, " : Transportando")
			godes.Advance(tim_gen.Get(5, 10))
			tr.state = TRANSPORTANDO
		case 'D':
			fmt.Println("T ", tr.id, " : Descargando")
			pils.lis.Front().Value.(stockPile).qe.Place(tr)
			pils.lis.Front().Value.(stockPile).empty_qe.Set(false)
			pils.nextList()
			tr.busy.Set(true)
			tr.busy.Wait(false)
			tr.state = DESCARGANDO
		case 'M':
			fmt.Println("T ", tr.id, " : Volviendo a cargar")
			godes.Advance(tim_gen.Get(1, 5))
			tr.state = REGRESANDO
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

func main() {
	pils = pils.init(STOCK_PILE_TYPE)
	chars = chars.init(CHARGER_TYPE)
	godes.AddRunner(&truck{&godes.Runner{}, 1, 0, godes.NewBooleanControl(), PARADO})
	godes.AddRunner(&truck{&godes.Runner{}, 2, 0, godes.NewBooleanControl(), PARADO})
	for i := 0; i < 3; i++ {
		pils.addList(stockPile{&godes.Runner{}, i, godes.NewFIFOQueue("pile"), godes.NewBooleanControl()})
	}
	for i := 0; i < 3; i++ {
		chars.addList(shovel{&godes.Runner{}, i, godes.NewFIFOQueue("chars"), godes.NewBooleanControl()})
	}
	godes.Run()
	godes.WaitUntilDone()
}
