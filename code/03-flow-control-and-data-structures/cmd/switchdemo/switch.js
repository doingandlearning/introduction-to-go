
const day = 2
let dayWord = ""

switch (day) {
	case 1:
		dayWord = "Monday"
		break
	case 2:
		dayWord = "Tuesday"
		break
	case 3:
		dayWord = "Wednesday"
		break
	case 4:
		dayWord = "Thursday"
		break
	case 5:
		dayWord = "Friday"
		break
	case 6: // Saturday
		dayWord = "Saturday"
		break
	case 7:
		dayWord = "Sunday"
		break
	default:
		dayWord = "Invalid day"
		break
}
