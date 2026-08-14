// Command fmtdemo exercises every fmt verb and function covered in the
// Topic 2 lecture, against one real struct, so you can match output
// lines back to the verb that produced them.
//
//	go run ./cmd/fmtdemo
package main

import (
	"bytes"
	"fmt"
	"os"
)

// Point is the struct used throughout this demo. It implements Stringer
// (below), so %v and Println use String() instead of the default
// struct format once that's set up.
type Point struct {
	X, Y int
}

// String implements fmt.Stringer. Once this exists, any %v or Println
// call involving a Point uses it automatically.
func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

// plainPoint has no String method, so it falls back to fmt's default
// struct formatting — used to show what %v looks like before Stringer
// enters the picture.
type plainPoint struct {
	X, Y int
}

func main() {
	fmt.Println("=== basic printers ===")
	fmt.Print("Print: ", "no trailing newline — ")
	fmt.Println("Println adds one")
	fmt.Printf("Printf: %s scores %d\n", "Go", 10)

	fmt.Println("\n=== core verbs against plainPoint (no Stringer) ===")
	pp := plainPoint{X: 3, Y: 4}
	fmt.Printf("%%v   -> %v\n", pp)
	fmt.Printf("%%+v  -> %+v\n", pp)
	fmt.Printf("%%#v  -> %#v\n", pp)
	fmt.Printf("%%T   -> %T\n", pp)

	fmt.Println("\n=== scalar verbs ===")
	fmt.Printf("%%d       -> %d\n", 42)
	fmt.Printf("%%s       -> %s\n", "hello")
	fmt.Printf("%%q       -> %q\n", "hello")
	fmt.Printf("%%f       -> %f\n", 3.14159)
	fmt.Printf("%%.2f     -> %.2f\n", 3.14159)
	fmt.Printf("%%t       -> %t\n", true)
	fmt.Printf("%%p       -> %p\n", &pp)

	fmt.Println("\n=== Sprintf and Errorf ===")
	msg := fmt.Sprintf("built as a string: %d items", 7)
	fmt.Println(msg)

	inner := fmt.Errorf("file not found")
	wrapped := fmt.Errorf("loading config: %w", inner)
	fmt.Println("wrapped error:", wrapped)

	fmt.Println("\n=== Fprintf / Fprintln against different writers ===")
	fmt.Fprintln(os.Stderr, "this line goes to stderr, not stdout")

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d bytes written into an in-memory buffer", 123)
	fmt.Println("read back from the buffer:", buf.String())

	fmt.Println("\n=== Stringer: Point controls its own %v ===")
	p := Point{X: 3, Y: 4}
	fmt.Printf("%%v on Point (has String())    -> %v\n", p)
	fmt.Println("Println on Point (has String) ->", p)
	// %+v also defers to String() once it exists — the + flag only adds
	// field names to the *default* struct format, and Stringer replaces
	// that default entirely, for every v-family verb.
	fmt.Printf("%%+v also uses String() here    -> %+v\n", p)
	// %#v is the exception: it always shows the Go-syntax literal,
	// ignoring Stringer, because its whole purpose is showing you the
	// real underlying value.
	fmt.Printf("%%#v ignores Stringer           -> %#v\n", p)
}
