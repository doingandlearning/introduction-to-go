// Command slicealias walks through the single most common new-Go-dev bug,
// one step at a time: slicing shares the underlying array, so mutating
// through a "view" mutates the original — until an append forces a
// reallocation and the sharing silently breaks.
//
//	go run ./cmd/slicealias
//
// Watch the len()/cap() printed at each step. The moment cap() jumps to a
// bigger number than it was, that's the reallocation — aliasing with
// original stops from that point forward.
package main

import "fmt"

func main() {
	original := []int{1, 2, 3, 4, 5}
	fmt.Println("step 0: original       =", original, "len", len(original), "cap", cap(original))

	// A sub-slice does NOT copy elements 1 and 2. It's a new small struct
	// (pointer + length + capacity) pointing at original's backing array.
	view := original[1:3]
	fmt.Println("step 1: view           =", view, "len", len(view), "cap", cap(view))

	// Mutate through the view. Nobody wrote to original directly.
	view[0] = 99
	fmt.Println("step 2: view[0] = 99")
	fmt.Println("        view           =", view)
	fmt.Println("        original       =", original, "<- changed, through view")

	// Append while there's still room in the shared backing array
	// (cap(view) is 4, len(view) is 2, so there's room for 2 more).
	// This OVERWRITES original[3] without original ever being touched
	// directly.
	view = append(view, 100)
	fmt.Println("step 3: view = append(view, 100)")
	fmt.Println("        view           =", view, "len", len(view), "cap", cap(view))
	fmt.Println("        original       =", original, "<- original[3] overwritten")

	// Append again, past capacity this time. Go reallocates a new,
	// larger backing array and copies the data over. From here on, view
	// and original are independent.
	view = append(view, 200, 300, 400)
	fmt.Println("step 4: view = append(view, 200, 300, 400)")
	fmt.Println("        view           =", view, "len", len(view), "cap", cap(view), "<- cap jumped: reallocated")

	view[0] = -1
	fmt.Println("step 5: view[0] = -1")
	fmt.Println("        view           =", view)
	fmt.Println("        original       =", original, "<- UNCHANGED: sharing broke at step 4")
}
