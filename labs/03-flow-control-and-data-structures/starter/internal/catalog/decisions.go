package catalog

// LateFeeTier classifies how overdue a book is.
//
// TODO(Exercise 2): implement using a switch with NO condition —
// switch { case daysLate >= 30: ... }
//
// Return "suspended" for 30+ days late, "warning" for 7-29 days,
// "none" for under 7 days.
func LateFeeTier(daysLate int) string {
	// TODO: replace with a conditionless switch.
	_ = daysLate
	return ""
}

// DeskSchedule reports the desk's hours for a given day, where day is
// 1 (Monday) through 7 (Sunday).
//
// TODO(Exercise 2): implement using switch day { ... }. Group
// Monday-Friday into one case returning "open 9am-6pm". Give Saturday
// its own case that falls through into Sunday's case deliberately, so
// both return "open 10am-2pm only".
func DeskSchedule(day int) string {
	// TODO: replace with a switch on day, using fallthrough from
	// Saturday's case into Sunday's case.
	_ = day
	return ""
}
