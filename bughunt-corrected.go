// bughunt_solution.go
// Fully corrected version of the practice assignment.
// All common Go errors have been fixed.
// Expected output (letter order may vary):
//
// JSON person: {"Name":"Alice","Age":30}
// Converted value: 42
// Slice after modify: [99 2 3]
// Counter after increments: 5
// Map after modify: map[a:99]
// Safe type assertion: not a string
// All values printed: a b c   (in some order)

package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

// 1. Visibility / JSON – fields must be exported (capitalized)
type person struct {
	Name string `json:"Name"` // was lowercase "name"
	Age  int    `json:"Age"`
}

// 2. Type mismatch – explicit conversion required
func convert(i int) float64 {
	return float64(i) // was: return i
}

// 3. Slice modification – works because the slice header points to the same underlying array
func modifySlice(s []int) {
	s[0] = 99
}

// 4. Struct method – pointer receiver so the original value is modified
type counter struct {
	n int
}

func (c *counter) increment() { // was value receiver (c counter)
	c.n++
}

// 5. Map modification – maps are reference types
func modifyMap(m map[string]int) {
	m["a"] = 99
}

// 6. Nil interface trap – never return a typed nil
type MyError struct {
	Msg string
}

func (e *MyError) Error() string {
	return e.Msg
}

func mayFail(shouldFail bool) error {
	if shouldFail {
		return &MyError{Msg: "something went wrong"}
	}
	return nil // explicit nil interface (was: return typed-nil *MyError)
}

// 7. Type assertion – always use the comma-ok form
func safeAssert(i interface{}) {
	if s, ok := i.(string); ok {
		fmt.Println("String value:", s)
	} else {
		fmt.Println("Safe type assertion: not a string")
	}
}

// 8. Loop variable + WaitGroup – fixed for correctness on all Go versions
func printValues(values []string) {
	var wg sync.WaitGroup
	for _, v := range values {
		wg.Add(1)
		v := v // capture per-iteration value (classic pre-1.22 fix; harmless on 1.22+)
		go func() {
			defer wg.Done() // was missing
			fmt.Print(v, " ")
		}()
	}
	wg.Wait()
	fmt.Println()
}

func main() {
	// ----- All previous compile-time errors removed -----
	// (unused import, unused variable, missing commas, := redeclaration)

	// 1. JSON
	p := person{Name: "Alice", Age: 30}
	b, err := json.Marshal(p)
	if err != nil {
		fmt.Println("JSON error:", err)
		return
	}
	fmt.Println("JSON person:", string(b))

	// 2. Type conversion
	fmt.Println("Converted value:", convert(42))

	// 3. Slice
	s := []int{1, 2, 3}
	modifySlice(s)
	fmt.Println("Slice after modify:", s)

	// 4. Counter with pointer receiver
	c := counter{n: 0}
	for i := 0; i < 5; i++ {
		c.increment()
	}
	fmt.Println("Counter after increments:", c.n)

	// 5. Map
	m := map[string]int{"a": 1}
	modifyMap(m)
	fmt.Println("Map after modify:", m)

	// 6. Nil interface trap – now correctly returns nil
	err = mayFail(false)
	if err != nil {
		fmt.Println("Unexpected error:", err) // this branch is no longer taken
	}

	// 7. Safe type assertion
	safeAssert(42)

	// 8. Concurrent printing
	values := []string{"a", "b", "c"}
	fmt.Print("All values printed: ")
	printValues(values)
}
