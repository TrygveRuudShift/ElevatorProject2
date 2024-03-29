package prioritize

import (
	"elevators/hardware"
)

func MotorActionOnDecisionDeadline(
	prioritizedDirection hardware.MotorDirection,
	currentOrders OrderSummary) hardware.MotorDirection {

	switch prioritizedDirection {
	case hardware.MD_Up:
		if currentOrders.AboveFloor {
			return hardware.MD_Up
		}

	case hardware.MD_Down:
		if currentOrders.BelowFloor {
			return hardware.MD_Down
		}

	case hardware.MD_Stop:
		return hardware.MD_Stop
	default:
		panic("Invalid recent direction on door close")
	}

	if currentOrders.AboveFloor {
		return hardware.MD_Up
	} else if currentOrders.BelowFloor {
		return hardware.MD_Down
	} else {
		return hardware.MD_Stop
	}
}

func MotorActionOnFloorArrival(
	prioritizedDirection hardware.MotorDirection,
	currentOrders OrderSummary) hardware.MotorDirection {

	if currentOrders.CabAtFloor {
		return hardware.MD_Stop
	}
	switch prioritizedDirection {
	case hardware.MD_Up:
		if currentOrders.UpAtFloor {
			return hardware.MD_Stop
		}
		if currentOrders.AboveFloor {
			return hardware.MD_Up
		}
		if currentOrders.DownAtFloor {
			return hardware.MD_Stop
		}
		if currentOrders.BelowFloor {
			return hardware.MD_Down
		}

	case hardware.MD_Down:
		if currentOrders.DownAtFloor {
			return hardware.MD_Stop
		}
		if currentOrders.BelowFloor {
			return hardware.MD_Down
		}
		if currentOrders.UpAtFloor {
			return hardware.MD_Stop
		}
		if currentOrders.AboveFloor {
			return hardware.MD_Up
		}

	case hardware.MD_Stop:
		return hardware.MD_Stop
	default:
		panic("Invalid recent direction on floor arrival")
	}
	return hardware.MD_Stop
}
