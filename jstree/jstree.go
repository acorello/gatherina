package jstree

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dop251/goja/ast"
	jsparser "github.com/dop251/goja/parser"
)

func MustParse(js io.Reader) *ast.Program {
	program, err := jsparser.ParseFile(nil, "", js, 0)
	if err != nil {
		log.Panic("Error parsing JS", err)
	}
	return program
}

func Walk(node ast.Node, feed func(ast.Node, int) bool) {
	walk(node, 0, feed)
}

func walk(node ast.Node, level int, feed func(ast.Node, int) bool) {
	// we don't feed wrapper nodes as they have no value nor structure
	switch node := node.(type) {
	case *ast.ExpressionStatement:
		walk(node.Expression, level, feed)
		return
	}
	_stepOn := func(node ast.Node) {
		walk(node, level+1, feed)
	}
	more := feed(node, level)
	if !more {
		return
	}
	switch node := node.(type) {
	case *ast.Identifier, *ast.StringLiteral:
		return // leaves
	case *ast.Program:
		for _, declaration := range node.DeclarationList {
			_stepOn(declaration)
		}
		for _, statement := range node.Body {
			_stepOn(statement)
		}
	case *ast.DotExpression:
		_stepOn(node.Left)
		_stepOn(&node.Identifier)
	case *ast.CallExpression:
		_stepOn(node.Callee)
		for _, expression := range node.ArgumentList {
			_stepOn(expression)
		}
	default:
		fmt.Fprintf(os.Stderr, "🛑 %T\n", node)
	}
}
