package pretty

import (
	"math"
	"reflect"
)

const defaultPrecision = 0.0000001

type customDiff struct {
}

type Comparator interface {
	Diff(a, b interface{}) (desc []string, ok bool)
	StructuredDiff(a, b interface{}) (desc []StructuredDiff, ok bool)
}

type Equals func(a, b interface{}) bool
type Float64Equals func(a, b float64) bool

type Options struct {
	customComparators        map[reflect.Type]Equals
	numericComparator        Float64Equals
	ignoreTypeNameDifference bool
	labelFieldNames          []string
}

func WithIgnoreTypeNameDiffs(ignore bool) func(*Options) {
	return func(s *Options) {
		s.ignoreTypeNameDifference = ignore
	}
}

func WithPrecision(precision float64) func(*Options) {
	return func(s *Options) {
		globalPrecision = precision
	}
}

func WithCustomComparators(customComparators map[reflect.Type]Equals) func(*Options) {
	return func(s *Options) {
		s.customComparators = customComparators
	}
}

var globalPrecision = defaultPrecision

func newMustAbsoluteDeltaLessThan(e float64) func(a, b float64) bool {
	return func(a, b float64) bool {
		return math.Abs(a-b) <= e+globalPrecision
	}
}

func WithLabelFields(labelFieldNames ...string) func(*Options) {
	return func(s *Options) {
		s.labelFieldNames = labelFieldNames
	}
}

// WithNumericEpsilon - sets the maximum tolerance of absolute difference of all numeric types
func WithNumericEpsilon(epsilon float64) func(*Options) {
	return func(s *Options) {
		s.numericComparator = newMustAbsoluteDeltaLessThan(epsilon)
	}
}

func NewCustomDiff(options ...func(*Options)) Comparator {
	opts := Options{
		customComparators: make(map[reflect.Type]Equals),
	}

	for _, o := range options {
		o(&opts)
	}
	return &customDiffPrinter{
		customComparators:        opts.customComparators,
		numericComparator:        opts.numericComparator,
		ignoreTypeNameDifference: opts.ignoreTypeNameDifference,
		labelNames:               opts.labelFieldNames,
	}
}

type customDiffPrinter struct {
	customComparators        map[reflect.Type]Equals
	numericComparator        Float64Equals
	ignoreTypeNameDifference bool
	labelNames               []string
}

func (c customDiffPrinter) Diff(a, b interface{}) (desc []string, ok bool) {
	diffPrinter{
		w:                        (*sbuf)(&desc),
		ignoreTypeNameDifference: c.ignoreTypeNameDifference,
		customComparators:        c.customComparators,
		numericComparator:        c.numericComparator,
		labels:                   NewLabels(c.labelNames...),
		aVisited:                 make(map[visit]visit),
		bVisited:                 make(map[visit]visit),
	}.diff(reflect.ValueOf(a), reflect.ValueOf(b))
	return desc, len(desc) == 0
}

func (c customDiffPrinter) StructuredDiff(a, b interface{}) (desc []StructuredDiff, ok bool) {
	descStr := make([]string, 0)
	structuredOut := NewStructuredDiffer()
	diffPrinter{
		w:                        (*sbuf)(&descStr),
		structuredOutput:         structuredOut,
		ignoreTypeNameDifference: c.ignoreTypeNameDifference,
		customComparators:        c.customComparators,
		numericComparator:        c.numericComparator,
		labels:                   NewLabels(c.labelNames...),
		aVisited:                 make(map[visit]visit),
		bVisited:                 make(map[visit]visit),
	}.diff(reflect.ValueOf(a), reflect.ValueOf(b))
	return structuredOut.Results(), len(structuredOut.Results()) == 0
}
