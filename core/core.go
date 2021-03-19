// Package core defines the core contracts of slurp. Core also defines primitives
// that work directly with the core types.
package core

import (
	"errors"
)

var (
	// ErrNotFound is returned by Env when a binding is not found
	// for a given symbol/name.
	ErrNotFound = errors.New("not found")

	// ErrInvalidName is returned by Env when the bind name is
	// invalid.
	ErrInvalidName = errors.New("invalid bind name")
)

// Env represents the environment in which forms are evaluated.
type Env interface {
	// Name returns the name of this env frame.
	Name() string

	// Namespace returns the active namespace.
	Namespace() Namespace

	// Parent returns the parent/outer env of this env. Returns nil
	// if this env is the root.
	Parent() Env

	// Scope of the current environment.
	Scope() Scope

	// Child returns a new env with given frame name and vars bound.
	// Returned env will have this env as parent/outer.
	Child(name string, vars map[string]Any) Env

	// WithNamespace returns the named environment, creating it if it does
	// not already exist.
	WithNamespace(Namespace) Env
}

// Scope is a named collection of variable bindings.
type Scope interface {
	// Bind creates a local binding with given name and value.
	Bind(name Symbol, val Any) error

	// Resolve resolves the symbol in this env and return its value
	// if found. Returns ErrNotFound if name is not found in this
	// env frame.
	Resolve(name Symbol) (Any, error)
}

// Namespace is a string that defaults to "main" when empty.
type Namespace string

// String returns the namespace, substituting the empty string
// for "main".
func (ns Namespace) String() string {
	if ns == "" {
		return "main"
	}

	return string(ns)
}

// Analyzer implementation is responsible for performing syntax analysis
// on given form.
type Analyzer interface {
	// Analyze performs syntax analysis and macro expansions if supported
	// to produce Expr values which can be valuated against Env.
	Analyze(env Env, form Any) (Expr, error)
}

// Expr represents an expression that can be evaluated against an env.
type Expr interface {
	// Eval executes the expr against the env and returns the results.
	// It can have side-effects on env. (e.g., DefExpr).
	Eval(env Env) (Any, error)
}

// Root traverses the env hierarchy until it reaches the root env
// and returns it.
func Root(env Env) Env {
	for env.Parent() != nil {
		env = env.Parent()
	}
	return env
}

// Eval performs syntax analysis of the given form to produce an Expr
// and evaluates the Expr against the given Env.
func Eval(env Env, analyzer Analyzer, form Any) (Any, error) {
	if analyzer == nil {
		analyzer = constAnalyzer{}
	}

	expr, err := analyzer.Analyze(env, form)
	if err != nil || expr == nil {
		return nil, err
	}
	return expr.Eval(env)
}

type constAnalyzer struct{}

func (constAnalyzer) Analyze(_ Env, form Any) (Expr, error) {
	return constExpr{Const: form}, nil
}

type constExpr struct{ Const Any }

func (ce constExpr) Eval(_ Env) (Any, error) { return ce.Const, nil }
