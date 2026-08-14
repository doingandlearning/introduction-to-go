package main

// farewell is unexported and lives in the same package as main.go, in a
// separate file. No import needed — package membership, not the file
// boundary, is what Go cares about.
func farewell(name string) string {
	return "Goodbye, " + name + "."
}
