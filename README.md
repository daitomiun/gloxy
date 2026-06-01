# Gloxy: The lox programming language in go
The go implementation of the `Lox` language

## Build the compiler

```bash
go run .
```

Or compile it's binary

```bash
go build .
```

## Regenerate statements and expressions

To generate the `expr.go` and `stmt.go` files under the `ast` package: 

1. Define the new expression under `too/generate/main.go`

2. Run the tool
```bash
go run ./tool/generate/main.go ./ast
```

## Why a Language interpreter?

Because it's fun. duh...

I wanted to learn how interpreters and compilers works so this is the way to do it! 

## Motivation
