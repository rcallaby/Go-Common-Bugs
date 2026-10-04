// bughunt.go
// Practice assignment: Find and fix the common Go errors.
// Expected final output (order of some lines may vary because of concurrency):
//
// JSON person: {"Name":"Alice","Age":30}
// Converted value: 42.0
// Slice after modify: [99 2 3]
// Counter after increments: 5
// Map after modify: map[a:99]
// Safe type assertion: not a string
// All values printed: a b c   (in some order)
// No panics, no deadlocks, no data races.

package main

import (
    "encoding/json"
    "fmt"
    "sync"
    "time"
)

// 1. Visibility / JSON
type person struct {
    name string // bug: should be exported
    Age  int
}

// 2. Type mismatch helper
func convert(i int) float64 {
    return i // bug: missing conversion
}

// 3. Slice modification
func modifySlice(s []int) {
    s[0] = 99
}

// 4. Struct method – accidental copy
type counter struct {
    n int
}

func (c counter) increment() { // bug: value receiver
    c.n++
}

// 5. Map modification
func modifyMap(m map[string]int) {
    m["a"] = 99
}

// 6. Nil interface trap
type MyError struct {
    Msg string
}

func (e *MyError) Error() string {
    return e.Msg
}

func mayFail(shouldFail bool) error {
    var err *MyError
    if shouldFail {
        err = &MyError{Msg: "something went wrong"}
    }
    return err // bug: typed nil
}

// 7. Type assertion
func safeAssert(i interface{}) {
    s := i.(string) // bug: unsafe assertion
    fmt.Println("String value:", s)
}

// 8. Loop variable capture + data race + possible leak
func printValues(values []string) {
    var wg sync.WaitGroup
    for _, v := range values {
        wg.Add(1)
        go func() {
            // bug: captures loop variable (pre-1.22 style) + missing Done
            fmt.Print(v, " ")
        }()
    }
    wg.Wait()
    fmt.Println()
}

func main() {
    // --- Compile-time & simple runtime bugs ---

    // Unused import / variable (will prevent compilation)
    import "os" // remove or use
    x := 42     // declared and not used

    // Missing trailing comma
    nums := []int{
        1,
        2
        3,
    }

    // Short declaration reuse
    err := error(nil)
    err := fmt.Errorf("redeclare") // no new variables on left side of :=

    // --- Behavioral bugs ---

    p := person{name: "Alice", Age: 30}
    b, _ := json.Marshal(p)
    fmt.Println("JSON person:", string(b)) // Age missing, Name missing

    fmt.Println("Converted value:", convert(42))

    s := []int{1, 2, 3}
    modifySlice(s)
    fmt.Println("Slice after modify:", s)

    c := counter{n: 0}
    for i := 0; i < 5; i++ {
        c.increment()
    }
    fmt.Println("Counter after increments:", c.n) // still 0

    m := map[string]int{"a": 1} // ok, but later we will force a nil case
    modifyMap(m)
    fmt.Println("Map after modify:", m)

    // Nil pointer (uncomment to see the panic)
    // var ptr *int
    // fmt.Println(*ptr)

    // Nil interface trap
    err = mayFail(false)
    if err != nil {
        fmt.Println("Unexpected error:", err) // this prints even on success
    }

    // Unsafe type assertion
    safeAssert(42)

    // Concurrency bugs
    values := []string{"a", "b", "c"}
    printValues(values)

    // Extra challenges (optional – introduce them one by one):
    // 1. Create a nil map and try to write to it → panic
    // 2. Send on a closed channel
    // 3. Cause a real deadlock with an unbuffered channel
    // 4. Introduce a data race on a shared counter and fix it with Mutex or atomic
}
