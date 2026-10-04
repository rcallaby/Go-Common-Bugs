# Go Error Hunting Practice Assignment

This repository provides a focused practice exercise for identifying and correcting many of the most common errors encountered when writing Go programs. The material is intended for developers who already know the basics of the language and wish to strengthen their ability to diagnose compile-time, runtime, and concurrency-related problems.

## Overview

Two files are supplied:

- `bughunt.go` – a deliberately broken program containing a representative selection of frequent Go mistakes.
- `bughunt_solution.go` – the fully corrected version of the same program, accompanied by brief explanatory comments.

Working through the broken file and then comparing it with the solution is an effective way to internalize the patterns that produce these errors and the idiomatic ways to resolve them.

## Files and Their Purpose

| File                    | Purpose                                                                 |
|-------------------------|-------------------------------------------------------------------------|
| `bughunt.go`            | Starting point. Contains intentional defects for you to discover and repair. |
| `bughunt_solution.go`   | Reference implementation. Shows one correct way to fix every issue.     |
| `README.md`             | This document.                                                          |

## How to Use the Assignment

### 1. Prepare the environment

Ensure a recent Go toolchain (1.21 or later recommended) is installed and available on your `PATH`.

```bash
go version
```

### 2. Begin with the broken program

Copy or open `bughunt.go`. Attempt to build and run it:

```bash
go build bughunt.go
go run bughunt.go
```

The compiler and runtime will surface a series of errors. Address them systematically:

- Resolve all compile-time issues first (unused imports or variables, missing commas, short-declaration conflicts, type mismatches, visibility problems, etc.).
- Eliminate runtime panics (nil-pointer dereferences, unsafe type assertions, sends on closed channels, etc.).
- Correct logical and concurrency defects (nil-interface trap, value versus pointer receivers, loop-variable capture, missing `WaitGroup` coordination, data races).

You may find it helpful to enable the race detector while iterating:

```bash
go run -race bughunt.go
```

### 3. Verify correctness

When the program compiles and runs cleanly, the output should resemble the following (the order of the concurrent letters may vary):

```
JSON person: {"Name":"Alice","Age":30}
Converted value: 42
Slice after modify: [99 2 3]
Counter after increments: 5
Map after modify: map[a:99]
Safe type assertion: not a string
All values printed: a b c
```

### 4. Compare with the solution

After you have produced a working version, open `bughunt_solution.go`. Review the comments that accompany each correction. This step reinforces the reasoning behind the fixes and highlights idiomatic Go practices.

You may run the solution directly for confirmation:

```bash
go run bughunt_solution.go
go run -race bughunt_solution.go
```

## Learning Objectives

By completing the exercise you will practice diagnosing and repairing:

- Unused imports and variables
- Exported versus unexported identifiers and their effect on JSON marshaling
- Absence of implicit type conversions
- Trailing-comma requirements in composite literals
- Restrictions on the short declaration operator (`:=`)
- Nil pointer dereferences
- Slice and map reference semantics
- Value versus pointer receivers
- The nil-interface trap
- Unsafe type assertions
- Common concurrency pitfalls (loop-variable capture, missing synchronization, data races)

## Optional Extensions

Once the core assignment is complete, consider extending the program with additional challenges such as:

- Protecting a shared counter with a mutex or atomic operations under high concurrency
- Demonstrating and then eliminating a goroutine leak
- Safely handling closed channels
- Introducing and resolving a deliberate deadlock

These extensions are left as further practice; the supplied solution does not implement them.

## License and Attribution

The code and accompanying documentation are provided for educational purposes. You are free to adapt them for personal study, classroom use, or internal training.

Thank you for working through this exercise. Careful attention to the errors illustrated here will help you write more robust and idiomatic Go code.

