package state

type State string

var (
	Auth    State = "Auth"
	Capture State = "Capture"
	Refund  State = "Refund"
	Void    State = "Void"
)

func TransitionAllowed(current State, targetState State) bool {
	switch current {
	case Auth:
		if targetState == Capture || targetState == Void {
			return true
		}
	case Capture:
		if targetState == Capture || targetState == Refund {
			return true
		}
	case Refund:
		if targetState == Refund {
			return true
		}
	case Void:
		return false
	}
	return false
}

// func HandleState(currentState State, request store.Transaction) {

// }

// func doAuth() {

// }

// func doCapture() {

// }

// func doRefund() {

// }

// func doVoid() {

// }
