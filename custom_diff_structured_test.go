package pretty

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_customStructuredDiffPrinter_Diff(t *testing.T) {
	type fields struct {
		opts []func(*Options)
	}
	type args struct {
		a interface{}
		b interface{}
	}

	type testStruct2 struct {
		str string
	}
	type testStruct3sameAstestStruct struct {
		intField   int
		floatField float64
		child      testStruct2
	}

	type testStruct struct {
		intField   int
		floatField float64
		child      testStruct2
	}

	type testTimeStruct struct {
		TimeField  time.Time
		timeField  time.Time
		floatField float64
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantDesc []StructuredDiff
		wantOk   bool
	}{
		{
			name: "equals",
			fields: fields{
				opts: nil,
			},
			args: args{
				a: testStruct{
					intField:   1,
					floatField: 2.3,
					child: testStruct2{
						str: "strValue",
					},
				},
				b: testStruct{
					intField:   1,
					floatField: 2.3,
					child: testStruct2{
						str: "strValue",
					},
				},
			},
			wantDesc: nil,
			wantOk:   true,
		},
		{
			name: "not equals",
			fields: fields{
				opts: nil,
			},
			args: args{
				a: testStruct{
					intField:   1,
					floatField: 2.3,
					child: testStruct2{
						str: "strValue",
					},
				},
				b: testStruct{
					intField:   2,
					floatField: 2.3,
					child: testStruct2{
						str: "strValue",
					},
				},
			},
			wantDesc: []StructuredDiff{
				{
					FieldName: "intField",
					Labels:    []Label{},
					ValueA:    "1",
					ValueB:    "2",
				},
			},
			wantOk: false,
		},
		{
			name: "numeric comparator",
			fields: fields{
				opts: []func(options *Options){
					WithNumericEpsilon(0.5),
				},
			},
			args: args{
				a: testStruct{
					intField:   1,
					floatField: 2.4,
					child: testStruct2{
						str: "strValue",
					},
				},
				b: testStruct{
					intField:   1,
					floatField: 2.3,
					child: testStruct2{
						str: "strValue",
					},
				},
			},
			wantDesc: nil,
			wantOk:   true,
		},
		{
			name: "numeric comparator 2",
			fields: fields{
				opts: []func(options *Options){
					WithNumericEpsilon(0.01),
				},
			},
			args: args{
				a: testStruct{
					intField:   1,
					floatField: 53.23,
					child: testStruct2{
						str: "strValue",
					},
				},
				b: testStruct{
					intField:   1,
					floatField: 53.24,
					child: testStruct2{
						str: "strValue",
					},
				},
			},
			wantDesc: nil,
			wantOk:   true,
		},
		{
			name: "type names ignore",
			fields: fields{
				opts: []func(options *Options){
					WithIgnoreTypeNameDiffs(true),
				},
			},
			args: args{
				a: testStruct3sameAstestStruct{
					intField:   1,
					floatField: 53.23,
					child: testStruct2{
						str: "strValue",
					},
				},
				b: testStruct{
					intField:   1,
					floatField: 53.23,
					child: testStruct2{
						str: "strValue",
					},
				},
			},
			wantDesc: nil,
			wantOk:   true,
		},
		{
			name: "label fields",
			fields: fields{
				opts: []func(options *Options){
					WithLabelFields("str"),
				},
			},
			args: args{
				a: testStruct{
					intField:   1,
					floatField: 53.23,
					child: testStruct2{
						str: "strValue A",
					},
				},
				b: testStruct{
					intField:   1,
					floatField: 53.23,
					child: testStruct2{
						str: "strValue B",
					},
				},
			},
			wantDesc: []StructuredDiff{
				{
					FieldName: "child.str",
					Labels:    []Label{{Name: "str", Value: "strValue A"}},
					ValueA:    "\"strValue A\"",
					ValueB:    "\"strValue B\"",
				},
			},
			wantOk: false,
		}, {
			name: "label fields 2",
			fields: fields{
				opts: []func(options *Options){
					WithLabelFields("str"),
				},
			},
			args: args{
				a: testStruct{
					intField:   1,
					floatField: 53.23,
					child: testStruct2{
						str: "strValue A",
					},
				},
				b: testStruct{
					intField:   1,
					floatField: 53.23,
					child: testStruct2{
						str: "strValue B",
					},
				},
			},
			wantDesc: []StructuredDiff{
				{
					FieldName: "child.str",
					Labels:    []Label{{Name: "str", Value: "strValue A"}},
					ValueA:    "\"strValue A\"",
					ValueB:    "\"strValue B\"",
				},
			},
			wantOk: false,
		},
		{
			name: "time fields",
			fields: fields{
				opts: []func(options *Options){},
			},
			args: args{
				a: testTimeStruct{
					TimeField:  time.Date(2022, time.July, 1, 1, 31, 31, 0, time.Local),
					timeField:  time.Date(2022, time.July, 1, 1, 31, 31, 0, time.Local),
					floatField: 53.23,
				},
				b: testTimeStruct{
					TimeField:  time.Date(2022, time.July, 2, 1, 31, 31, 0, time.Local),
					timeField:  time.Date(2022, time.July, 2, 1, 31, 31, 0, time.Local),
					floatField: 53.23,
				},
			},
			wantDesc: []StructuredDiff{
				{
					FieldName: "TimeField",
					Labels:    []Label{},
					ValueA:    "2022-07-01 01:31:31 +0200 CEST",
					ValueB:    "2022-07-02 01:31:31 +0200 CEST",
				},
				{
					FieldName: "timeField.ext",
					Labels:    []Label{},
					ValueA:    "63792228691",
					ValueB:    "63792315091",
				},
			},
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCustomDiff(tt.fields.opts...)
			gotDesc, gotOk := c.StructuredDiff(tt.args.a, tt.args.b)
			if !assert.Equal(t, tt.wantDesc, gotDesc) {
				t.Errorf("Diff() gotDesc = \n%#v, want \n%#v", gotDesc, tt.wantDesc)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Diff() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
