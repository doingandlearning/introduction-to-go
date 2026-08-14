// Command greetdemo imports the local "greeting" package to show
// exported vs. unexported identifiers crossing a package boundary.
//
// Try this: open internal/greeting/greeting.go, rename greetingPrefix to
// GreetingPrefix, and call greeting.GreetingPrefix() from here. It works.
// Now revert the rename but leave the call to GreetingPrefix() in place —
// the build fails, because you're reaching for something that no longer
// exists outside the package.
package main

import (
	"fmt"

	"example.com/go-essentials/internal/greeting"
)

func main() {
	fmt.Println(greeting.Greet("delegates"))
}
