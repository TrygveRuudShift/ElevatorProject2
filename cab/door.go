package cab

import (
	"elevators/eta"
	"elevators/hardware"
	"elevators/orders"
	"elevators/prioritize"
	"elevators/timer"
)

func setDoorAndCabState(state hardware.DoorAction) {
	switch state {
	case hardware.DS_Open_Cab:
		openDoor()
		orders.CompleteOrderCab(Cab.AboveOrAtFloor)

	case hardware.DS_Open_Up:
		openDoor()
		orders.CompleteOrderCabAndUp(Cab.AboveOrAtFloor)
		Cab.RecentDirection = hardware.MD_Up

	case hardware.DS_Open_Down:
		openDoor()
		orders.CompleteOrderCabAndDown(Cab.AboveOrAtFloor)
		Cab.RecentDirection = hardware.MD_Down

	case hardware.DS_Close:
		closeDoor()
		timer.DecisionDeadlineTimer.TimerStart()

	case hardware.DS_Do_Nothing:
		break
	}
}

func openDoor() {
	hardware.SetDoorOpenLamp(true)
	Cab.Behaviour = DoorOpen
	timer.DoorTimer.TimerStart()
}

func closeDoor() {
	hardware.SetDoorOpenLamp(false)
	Cab.Behaviour = Idle
	timer.DoorTimer.TimerStop()
}

func FSMObstructionChange(
	obstructed bool,
	allOrders orders.AllOrders) {

	Cab.DoorObstructed = obstructed

	switch obstructed {
	case true:
		switch Cab.Behaviour {
		case DoorOpen:
			Cab.Behaviour = CabObstructed
		}

	case false:
		switch Cab.Behaviour {
		case CabObstructed:
			Cab.Behaviour = DoorOpen
			if timer.DoorTimer.TimedOut() {
				FSMDoorTimeout(allOrders)
			}
		}
	}
}

func FSMDoorTimeout(allOrders orders.AllOrders) ElevatorBehaviour {

	switch Cab.Behaviour {
	case DoorOpen:
		currentOrderSummary := orders.GetOrderSummary(
			allOrders,
			Cab.AboveOrAtFloor)

		prioritizedDirection := eta.PrioritizedDirection(
			Cab.AboveOrAtFloor,
			Cab.RecentDirection,
			allOrders,
			eta.GetInternalETAs())

		doorAction := prioritize.DoorActionOnDoorTimeout(
			prioritizedDirection,
			Cab.DoorObstructed,
			currentOrderSummary)

		setDoorAndCabState(doorAction)

		if doorAction == hardware.DS_Close {
			eta.UpdateOrderAndInternalETAs(
				Cab.RecentDirection,
				Cab.AboveOrAtFloor,
				false,
				allOrders)

			timer.DecisionDeadlineTimer.TimerStart()
		}

	case CabObstructed:
		break

	default:
		panic("Invalid cab state on door timeout")
	}

	return Cab.Behaviour
}

func FSMFloorStop(
	floor int,
	allOrders orders.AllOrders) ElevatorBehaviour {

	switch Cab.Behaviour {
	case Idle:
		currentOrderSummary := orders.GetOrderSummary(
			allOrders,
			Cab.AboveOrAtFloor)

		prioritizedDirection := eta.PrioritizedDirection(
			Cab.AboveOrAtFloor,
			Cab.RecentDirection,
			allOrders,
			eta.GetInternalETAs())

		doorAction := prioritize.DoorActionOnFloorStop(
			prioritizedDirection,
			currentOrderSummary)

		setDoorAndCabState(doorAction)

	default:
		break
	}

	return Cab.Behaviour
}
