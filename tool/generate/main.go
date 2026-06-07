package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage: generate_ast <output directory>")
		os.Exit(64)
	}
	outputDir := args[0]
	defineAst(outputDir, "Expr", []string{
		"Ternary  : Condition Expr, Then Expr, Else Expr",
		"Assign   : Name Token, Value Expr",
		"Binary   : Left Expr, Operator Token, Right Expr",
		"Grouping : Expression Expr",
		"Literal  : Value Expr",
		"Unary    : Operator Token, Right Expr",
		"Variable : Name Token",
		"Logical  : Left Expr, Operator Token, Right Expr",
	})
	defineAst(outputDir, "Stmt", []string{
		"BlockStmt      : Statements []Stmt",
		"ExpressionStmt : Expression Expr",
		"PrintStmt      : Expression Expr",
		"VarStmt        : Name Token, Initializer Expr",
		"IfStmt					: Condition Expr, ThenBranch Stmt, ElseBranch Stmt",
		"WhileStmt      : Condition Expr, Body Stmt",
	})
}

func defineAst(outputDir, baseName string, types []string) {
	path := filepath.Join(outputDir, strings.ToLower(baseName)+".go")
	fmt.Printf("path -> %s\n", path)
	f, err := os.Create(path)
	check(err)
	defer f.Close()
	_, err = fmt.Fprintln(f, "package ast")
	check(err)
	_, err = fmt.Fprintln(f)
	check(err)
	// interface methods implemented here
	_, err = fmt.Fprintf(f, "type %s any\n", baseName)
	check(err)

	// The AST structs
	for _, t := range types {
		structName := strings.TrimSpace(strings.Split(t, ":")[0])
		fields := strings.TrimSpace(strings.Split(t, ":")[1])
		defineType(f, structName, fields)
	}

	fmt.Println("Done")
}

func defineType(w io.Writer, structName, fieldList string) {
	_, err := fmt.Fprintf(w, "type %s struct {\n", structName)
	check(err)

	for field := range strings.SplitSeq(fieldList, ", ") {
		_, err := fmt.Fprintf(w, "	%s\n", field)
		check(err)
	}
	_, err = fmt.Fprintln(w, "}")
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
