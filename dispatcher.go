package main
import "container/list"

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