package catalog

// LateFeeTier classifies how overdue a book is, using a switch with no
// condition — a clean if/else chain written as a switch.
func LateFeeTier(daysLate int) string {
	switch {
	case daysLate >= 30:
		return "suspended"
	case daysLate >= 7:
		return "warning"
	default:
		return "none"
	}
}

// DeskSchedule reports the desk's hours for a given day, where day is
// 1 (Monday) through 7 (Sunday). Saturday deliberately falls through
// into Sunday's case.
func DeskSchedule(day int) string {
	switch day {
	case 1, 2, 3, 4, 5:
		return "open 9am-6pm"
	case 6:
		fallthrough
	case 7:
		return "open 10am-2pm only"
	default:
		return "unknown day"
	}
}
