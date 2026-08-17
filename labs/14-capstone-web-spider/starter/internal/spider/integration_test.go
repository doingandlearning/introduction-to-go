package spider

import "testing"

// TODO (Exercise 5): write TestEndToEnd yourself — this is the first
// time since Topic 12 that you're authoring a test rather than making a
// pre-written one pass. Wire urlfile.ReadURLs -> Run -> WriteCSV
// together against a real httptest.Server (like fetch_test.go already
// does) and a real temp file (like urlfile_test.go already does), then
// assert the report file actually reflects what the fake server
// returned. See the exercise sheet for the exact shape it should take.
func TestEndToEnd(t *testing.T) {
	t.Skip("TODO: implement TestEndToEnd")
}
