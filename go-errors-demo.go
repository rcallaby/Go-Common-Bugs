// go_errors_demo.go
// Demonstrates common Go errors / gotchas.
// Many compile-time errors are shown as comments so the file still builds.
// Runtime demos are controlled by flags.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// 1. Compiler & Syntax Errors
// ---------------------------------------------------------------------------

// --- Unused Variables / Imports ---
// The compiler rejects unused imports and declared-but-unused variables.
//
// import "os"          // "imported and not used: "os""
//
// func unusedDemo() {
// 	x := 42             // "x declared and not used"
// 	_ = x               // silence with blank identifier
// }

// --- Capitalization = Visibility (exported vs unexported) ---
type person struct {
	Name string // exported (JSON will see it)
	age  int    // unexported (JSON will ignore it)
}

func capitalizationDemo() {
	p := person{Name: "Alice", age: 30}
	b, _ := json.Marshal(p)
	fmt.Println("JSON (age missing because unexported):", string(b))
	// Output: {"Name":"Alice"}
}

// --- Type mismatches (no implicit conversion) ---
func typeMismatchDemo() {
	var i int = 42
	// var f float64 = i          // cannot use i (type int) as type float64
	var f float64 = float64(i) // explicit conversion required
	fmt.Println("Converted:", f)
}

// --- Missing trailing commas in composite literals ---
func missingCommaDemo() {
	// This fails to compile:
	// s := []int{
	// 	1,
	// 	2          // syntax error: unexpected newline, expecting comma or }
	// }
	s := []int{
		1,
		2, // trailing comma is mandatory on multi-line literals
	}
	fmt.Println("Slice with trailing commas:", s)
}

// --- no new variables on left side of := ---
func shortDeclDemo() {
	err := error(nil)
	// later...
	// err := someFunc()   // no new variables on left side of :=
	err = fmt.Errorf("reassigned") // use = instead
	_ = err
}

// ---------------------------------------------------------------------------
// 2. Pointers, Memory & Structs
// ---------------------------------------------------------------------------

func nilPointerDemo() {
	var p *int // nil by default
	// fmt.Println(*p) // panic: runtime error: invalid memory address or nil pointer dereference
	if p != nil {
		fmt.Println(*p)
	} else {
		fmt.Println("Safely detected nil pointer")
	}
}

// Slices / maps are reference types (header is copied, underlying data shared)
func sliceMapDemo() {
	s := []int{1, 2, 3}
	modifySlice(s)
	fmt.Println("After modifySlice (shared underlying array):", s) // [99 2 3]

	s2 := []int{1, 2, 3}
	appendWrong(s2)
	fmt.Println("After appendWrong (capacity not exceeded, still shared):", s2)

	m := map[string]int{"a": 1}
	modifyMap(m)
	fmt.Println("After modifyMap:", m) // map[a:99]
}

func modifySlice(s []int)        { s[0] = 99 }
func appendWrong(s []int)        { _ = append(s, 4) } // caller does not see the append if capacity is exceeded
func modifyMap(m map[string]int) { m["a"] = 99 }

// Accidental struct copying (value vs pointer receiver)
type counter struct{ n int }

func (c counter) incrementValue() { c.n++ } // operates on a copy
func (c *counter) incrementPtr()  { c.n++ } // operates on the original

func structCopyDemo() {
	c := counter{n: 0}
	c.incrementValue()
	fmt.Println("After value receiver:", c.n) // still 0
	c.incrementPtr()
	fmt.Println("After pointer receiver:", c.n) // 1
}

// JSON into nil map/slice
func jsonNilDemo() {
	var m map[string]string // nil
	// json.Unmarshal([]byte(`{"k":"v"}`), &m) // works in modern Go (allocates),
	// but older mental models / some edge cases still expect make()
	m = make(map[string]string)
	_ = json.Unmarshal([]byte(`{"k":"v"}`), &m)
	fmt.Println("Unmarshaled map:", m)
}

// ---------------------------------------------------------------------------
// 3. Concurrency & Goroutines
// ---------------------------------------------------------------------------

func goroutineLeakDemo() {
	ch := make(chan int)
	go func() {
		// This goroutine will block forever (leak) if nothing ever receives
		ch <- 42
	}()
	// Intentionally never receive → leak
	fmt.Println("Goroutine launched that will leak (blocked on send)")
	time.Sleep(100 * time.Millisecond)
}

func deadlockDemo() {
	ch := make(chan int)
	// ch <- 1               // fatal error: all goroutines are asleep - deadlock!
	go func() { ch <- 1 }() // correct: have a receiver
	fmt.Println("Received:", <-ch)
}

func dataRaceDemo() {
	// Run with: go run -race go_errors_demo.go -race
	var counter int
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // data race!
		}()
	}
	wg.Wait()
	fmt.Println("Counter (racy, may be wrong):", counter)
}

func closedChannelDemo() {
	ch := make(chan int, 1)
	close(ch)
	// ch <- 1 // panic: send on closed channel
	fmt.Println("Channel already closed – sending would panic")
}

func unbufferedNoReceiverDemo() {
	ch := make(chan int)
	// ch <- 42 // blocks forever (or deadlock if only goroutine)
	go func() { ch <- 42 }()
	fmt.Println("Sent on unbuffered with a receiver ready:", <-ch)
}

// ---------------------------------------------------------------------------
// 4. Interfaces, Errors & Type System
// ---------------------------------------------------------------------------

// The famous nil-interface trap
type MyError struct{ Msg string }

func (e *MyError) Error() string { return e.Msg }

func mayFail(fail bool) error {
	var err *MyError // typed nil
	if fail {
		err = &MyError{Msg: "boom"}
	}
	return err // returns non-nil interface containing nil pointer!
}

func nilInterfaceDemo() {
	err := mayFail(false)
	fmt.Printf("err == nil? %v  (surprising!)\n", err == nil) // false
	fmt.Printf("type: %T  value: %v\n", err, err)

	// Correct pattern:
	if fail := false; fail {
		return
	}
	// return nil explicitly
}

// Ignoring errors
func ignoreErrorDemo() {
	// data, _ := os.ReadFile("nonexistent") // silent failure
	fmt.Println("Always check errors – never ignore with _ unless intentional")
}

// Type assertion panic
func typeAssertDemo() {
	var i interface{} = 42
	// s := i.(string) // panic: interface conversion: interface {} is int, not string
	if s, ok := i.(string); ok {
		fmt.Println(s)
	} else {
		fmt.Println("Safe type assertion failed as expected")
	}
}

// ---------------------------------------------------------------------------
// 5. Loop variable capture (pre-Go 1.22 behavior)
// ---------------------------------------------------------------------------
// In Go ≥ 1.22 this is fixed (per-iteration variables).
// On older compilers / go.mod < 1.22 the classic bug still appears.

func loopCaptureDemo() {
	values := []string{"a", "b", "c"}
	var wg sync.WaitGroup
	for _, v := range values {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("Captured (pre-1.22 would all be 'c'):", v)
		}()
	}
	wg.Wait()
}

// ---------------------------------------------------------------------------
// Main – select which demos to run
// ---------------------------------------------------------------------------

func main() {
	nilPtr := flag.Bool("nilptr", false, "demonstrate nil pointer dereference")
	nilIface := flag.Bool("niliface", false, "demonstrate nil interface trap")
	race := flag.Bool("race", false, "demonstrate data race (use with -race flag)")
	deadlock := flag.Bool("deadlock", false, "demonstrate deadlock")
	closed := flag.Bool("closedchan", false, "demonstrate send on closed channel")
	loop := flag.Bool("loopcapture", false, "demonstrate loop variable capture")
	flag.Parse()

	fmt.Println("=== Safe / always-run demos ===")
	capitalizationDemo()
	typeMismatchDemo()
	missingCommaDemo()
	shortDeclDemo()
	sliceMapDemo()
	structCopyDemo()
	jsonNilDemo()
	typeAssertDemo()
	ignoreErrorDemo()

	if *nilPtr {
		fmt.Println("\n=== Nil pointer demo ===")
		nilPointerDemo()
	}
	if *nilIface {
		fmt.Println("\n=== Nil interface trap ===")
		nilInterfaceDemo()
	}
	if *race {
		fmt.Println("\n=== Data race demo (run with go run -race) ===")
		dataRaceDemo()
	}
	if *deadlock {
		fmt.Println("\n=== Deadlock demo ===")
		deadlockDemo()
	}
	if *closed {
		fmt.Println("\n=== Closed channel demo ===")
		closedChannelDemo()
	}
	if *loop {
		fmt.Println("\n=== Loop variable capture ===")
		loopCaptureDemo()
	}

	// Intentionally leaky / blocking demos are left commented or behind flags
	// so the program exits cleanly by default.
	_ = goroutineLeakDemo // just to silence unused
}
