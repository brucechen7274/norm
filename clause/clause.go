package clause

import "errors"

var (
	// ErrInvalidClauseParams indicates that the argument to the clause is invalid
	ErrInvalidClauseParams = errors.New("invalid clause params")
)

// Interface clause interface
type Interface interface {
	// Name returns the name of the clause
	Name() string

	// MergeIn merges the same clause into the Clause object
	MergeIn(clause *Clause)

	// Build constructs the nGQL statement
	Build(nGQL Builder) error
}

// Clause is the general structure of a clause, containing a name and an expression
type Clause struct {
	Name       string
	Expression Expression
}

// Build  clause
func (c Clause) Build(nGQL Builder) error {
	return c.Expression.Build(nGQL)
}

// Options represents configuration options for a clause
type Options struct {
	PropNames []string // list of specified properties
	TagName   string   // specified tag name
}

type Option func(*Options)

// WithPropNames sets the specified property field names
func WithPropNames(propNames []string) Option {
	return func(o *Options) {
		o.PropNames = propNames
	}
}

// WithTagName sets the tag name used in the clause
func WithTagName(tagName string) Option {
	return func(o *Options) {
		o.TagName = tagName
	}
}
