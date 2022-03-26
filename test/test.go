package main

import (
	"elevators/controlunit/orderstate"
	"elevators/filesystem"
	"elevators/hardware"
	"fmt"
	"time"
)

// {
// filename := "test"
// filestate := orderstate.GetInternalETAs()
// if len(os.Args) > 1 {
// 	file, _ := json.MarshalIndent(filestate, "", " ")
// 	_ = ioutil.WriteFile("testresults/" + filename+os.Args[1]+".json", file, 0644)
// } else {
// 	file, _ := json.MarshalIndent(filestate, "", " ")
// 	_ = ioutil.WriteFile("testresults/" + filename+".json", file, 0644)
// }
// }

func hasOrder(inputState orderstate.OrderState) bool {
	//TODO: test (Possible riv ruskende here)
	return inputState.LastOrderTime.After(inputState.LastCompleteTime)
}

func testHasOrder() {
	fmt.Println("Testing hasOrder")
	fmt.Println("")

	// var time1 = time.Time{}
	// var time2 = time.Time{}
	orders := []orderstate.OrderState{
		orderstate.OrderState{},
		orderstate.OrderState{LastOrderTime: time.Time{}, LastCompleteTime: time.Time{}},
		orderstate.OrderState{LastOrderTime: time.Now(), LastCompleteTime: time.Time{}},
		orderstate.OrderState{LastOrderTime: time.Now(), LastCompleteTime: time.Now().Add(-1)},
		orderstate.OrderState{LastOrderTime: time.Now(), LastCompleteTime: time.Now().Add(-1 * time.Second)},
		orderstate.OrderState{LastOrderTime: time.Now(), LastCompleteTime: time.Now()}}
	for _, order := range orders {
		fmt.Print("Last order: ")
		fmt.Print(order.LastOrderTime)
		fmt.Print(", Last Complete: ")
		fmt.Print(order.LastCompleteTime)
		fmt.Print(" ; hasOrder : ")
		fmt.Println(hasOrder(order))
		fmt.Println("")
	}
}

func testAnyOrders() {
	filesystem.Init()
	orders := filesystem.ReadOrders()
	fmt.Print("Any orders: ")
	fmt.Println(orderstate.AnyOrders(orders))
}

func testFirstExternalETA() {
	filesystem.Init()
	orders := filesystem.ReadOrders()
	fmt.Print("first eta expire: ")
	fmt.Println(orderstate.FirstBestETAexpirationWithOrder(orders))
}

func testComputeETAs() {
	filesystem.Init()
	orders := filesystem.ReadOrders()
	cab := filesystem.ReadCabState()
	etas := orderstate.ComputeETAs(cab.AboveOrAtFloor,
		hardware.MD_Stop,
		cab.RecentDirection,
		orders)
	filesystem.Write("testresults/"+"computeETAs.json", etas)
	etas = orderstate.ComputeETAs(cab.AboveOrAtFloor,
		hardware.MD_Stop,
		hardware.MD_Down,
		orders)
	filesystem.Write("testresults/"+"computeETAsDown.json", etas)
	etas = orderstate.ComputeETAs(cab.AboveOrAtFloor,
		hardware.MD_Stop,
		hardware.MD_Up,
		orders)
	filesystem.Write("testresults/"+"computeETAsUp.json", etas)
}

func main() {
	// testHasOrder()
	testAnyOrders()
	testFirstExternalETA()
	testComputeETAs()
}
