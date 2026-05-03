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
		"Binary   : Left Expr, Operator Token, Right Expr",
		"Grouping : Expression Expr",
		"Literal  : Value any",
		"Unary    : Operator Token, Right Expr",
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
	_, err = fmt.Fprintf(f, "type %s interface{}\n", baseName)
	check(err)

	// The AST structs
	for _, t := range types {
		structName := strings.TrimSpace(strings.Split(t, ":")[0])
		fields := strings.TrimSpace(strings.Split(t, ":")[1])
		defineType(f, structName, fields)
	}

	// evaluator function
	defineEvaluator(f, baseName, types)

	fmt.Println("Done")
}

func defineEvaluator(w io.Writer, baseName string, types []string) {
	_, err := fmt.Fprintf(w, "func Evaluate(e %s) any {\n", baseName)

	_, err = fmt.Fprintln(w, "	switch t := e.(type) {")
	check(err)
	// case evaluators are implemented here
	for _, t := range types {
		typeName := strings.TrimSpace(strings.Split(t, ":")[0])
		_, err = fmt.Fprintf(w, "		case %s:\n", typeName)
		check(err)
		_, err = fmt.Fprintln(w, "			return t")
		check(err)
	}
	_, err = fmt.Fprintln(w, "		default:")
	check(err)
	_, err = fmt.Fprintln(w, `			return nil`)
	check(err)

	_, err = fmt.Fprintln(w, "	}")
	check(err)
	_, err = fmt.Fprintln(w, "}")
	check(err)
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
