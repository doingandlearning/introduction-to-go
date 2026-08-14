// Command wordcount builds a map[string]int word-frequency counter over a
// sentence, uses the comma-ok idiom to distinguish "appeared zero times"
// from "never seen at all," and then ranges over the same map twice to
// show that Go deliberately does not guarantee iteration order.
//
//	go run ./cmd/wordcount
package main

import (
	"fmt"
	"strings"
)

// countWords splits s on whitespace and counts occurrences of each word.
func countWords(s string) map[string]int {
	counts := make(map[string]int) // must initialize before writing
	for _, word := range strings.Fields(s) {
		word = strings.ToLower(strings.Trim(word, ".,!?"))
		counts[word]++
	}
	return counts
}

func main() {
	sentence := "the quick brown fox jumps over the lazy dog the fox runs"
	counts := countWords(sentence)

	fmt.Println("-- comma-ok: zero occurrences vs never seen --")
	for _, word := range []string{"fox", "the", "cat"} {
		n, seen := counts[word]
		if seen {
			fmt.Printf("%-6s appeared %d time(s)\n", word, n)
		} else {
			fmt.Printf("%-6s never appeared in the sentence\n", word)
		}
	}
	// Without comma-ok, counts["cat"] would just return 0 — indistinguishable
	// from a word that appeared zero times after being explicitly tracked.
	// `seen` is what tells them apart.

	fmt.Println()
	fmt.Println("-- range over the same map, twice --")
	fmt.Print("pass 1: ")
	for k := range counts {
		fmt.Print(k, " ")
	}
	fmt.Println()

	fmt.Print("pass 2: ")
	for k := range counts {
		fmt.Print(k, " ")
	}
	fmt.Println()
	fmt.Println("Run this program again (a fresh execution) and compare — the order is deliberately randomized, not stable across runs.")
}
