package util

import (
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type Underlying struct {
	Value  int
	Stamp  time.Time
	StampT *time.Time
	I      decimal.Decimal
}

type TypeA struct {
	A string
	U Underlying
}

type TypeB struct {
	B string
	U Underlying
}

func TestCompareData(t *testing.T) {
	var nilStruct *TypeA
	now := time.Now()
	tomorrow := now.Add(time.Hour * 24)

	type args struct {
		copy     interface{}
		original interface{}
		depth    uint
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "zero depth",
			args: args{
				copy:     &TypeA{A: "a"},
				original: nil,
				depth:    0,
			},
			wantErr: false,
		},
		{
			name: "interface copy/original is nil",
			args: args{
				copy:     nilStruct,
				original: nilStruct,
				depth:    1,
			},
			wantErr: false,
		},
		{
			name: "one of the type is nil",
			args: args{
				copy:     &TypeA{A: "a"},
				original: nilStruct,
				depth:    1,
			},
			wantErr: true,
		},
		{
			name: "different type",
			args: args{
				copy:     &TypeA{A: "a"},
				original: nil,
				depth:    1,
			},
			wantErr: true,
		},
		{
			name: "different types of interface copy and original",
			args: args{
				copy:     &TypeA{},
				original: &TypeB{},
				depth:    1,
			},
			wantErr: true,
		},
		{
			name: "success compare but got different value",
			args: args{
				copy:     &TypeA{A: "a"},
				original: &TypeA{A: "b"},
				depth:    1,
			},
			wantErr: true,
		},
		{
			name: "success compare",
			args: args{
				copy:     &TypeA{A: "a"},
				original: &TypeA{A: "a"},
				depth:    1,
			},
			wantErr: false,
		},
		{
			name: "success compare recursive",
			args: args{
				copy:     &TypeA{A: "a", U: Underlying{Value: 1, Stamp: now, StampT: &now, I: decimal.NewFromFloat(100)}},
				original: &TypeA{A: "a", U: Underlying{Value: 1, Stamp: now, StampT: &now, I: decimal.NewFromFloat(100)}},
				depth:    999,
			},
			wantErr: false,
		},
		{
			name: "success compare recursive - 2",
			args: args{
				copy:     &TypeA{A: "a", U: Underlying{Value: 1, Stamp: now, StampT: nil, I: decimal.NewFromFloat(100)}},
				original: &TypeA{A: "a", U: Underlying{Value: 1, Stamp: now, StampT: nil, I: decimal.NewFromFloat(100)}},
				depth:    999,
			},
			wantErr: false,
		},
		{
			name: "fail compare in recursive - Different nil values",
			args: args{
				copy:     &TypeA{A: "a", U: Underlying{Value: 1, Stamp: now, StampT: &now, I: decimal.NewFromFloat(100)}},
				original: &TypeA{A: "a", U: Underlying{Value: 1, Stamp: now, StampT: nil, I: decimal.NewFromFloat(100)}},
				depth:    999,
			},
			wantErr: true,
		},
		{
			name: "fail compare in recursive - Different time pointer value",
			args: args{
				copy:     &TypeA{A: "a", U: Underlying{Value: 1, Stamp: now, StampT: &now}},
				original: &TypeA{A: "a", U: Underlying{Value: 1, Stamp: tomorrow, StampT: &tomorrow}},
				depth:    999,
			},
			wantErr: true,
		},
		{
			name: "fail compare in recursive - Decimal",
			args: args{
				copy:     &TypeA{A: "a", U: Underlying{Value: 1, Stamp: now, StampT: &now, I: decimal.NewFromFloat(100)}},
				original: &TypeA{A: "a", U: Underlying{Value: 1, Stamp: now, StampT: &now, I: decimal.NewFromFloat(1000)}},
				depth:    999,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CompareData(tt.args.copy, tt.args.original, tt.args.depth); (err != nil) != tt.wantErr {
				t.Errorf("CompareData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCompareItems(t *testing.T) {
	type CustomStruct struct {
		A int
		B string
		C time.Time
		D decimal.Decimal
	}

	now := time.Now()
	testCases := []struct {
		name        string
		source      []CustomStruct
		target      []CustomStruct
		depth       uint
		wantErr     bool
		expectedErr string
	}{
		{
			name: "Depth is zero",
			source: []CustomStruct{
				{A: 1, B: "test", C: time.Now(), D: decimal.NewFromInt(1)},
			},
			target: []CustomStruct{
				{A: 1, B: "test", C: time.Now(), D: decimal.NewFromInt(1)},
			},
			depth:       0,
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "Different lengths",
			source: []CustomStruct{
				{A: 1, B: "test1"},
				{A: 2, B: "test2"},
			},
			target: []CustomStruct{
				{A: 1, B: "test1"},
			},
			depth:       1,
			wantErr:     true,
			expectedErr: "source and target has different length. Source: 2, Target: 1",
		},
		{
			name: "Matching slices with depth 1",
			source: []CustomStruct{
				{A: 1, B: "test", C: now, D: decimal.NewFromInt(10)},
			},
			target: []CustomStruct{
				{A: 1, B: "test", C: now, D: decimal.NewFromInt(10)},
			},
			depth:       1,
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "Mismatch at field A",
			source: []CustomStruct{
				{A: 1, B: "test"},
			},
			target: []CustomStruct{
				{A: 2, B: "test"},
			},
			depth:       1,
			wantErr:     true,
			expectedErr: "mismatch on array item index {0}. Details: field A is not matching. Source: 1, Target: 2",
		},
		{
			name: "Timestamp mismatch in field C",
			source: []CustomStruct{
				{A: 1, B: "test", C: now},
			},
			target: []CustomStruct{
				{A: 1, B: "test", C: now.Add(time.Hour)},
			},
			depth:       2,
			wantErr:     true,
			expectedErr: fmt.Sprintf("mismatch on array item index {0}. Details: field C is not matching. Source: %s, Target: %s", now, now.Add(time.Hour)),
		},
		{
			name: "Decimal mismatch in field D",
			source: []CustomStruct{
				{A: 1, B: "test", D: decimal.NewFromInt(10)},
			},
			target: []CustomStruct{
				{A: 1, B: "test", D: decimal.NewFromInt(20)},
			},
			depth:       2,
			wantErr:     true,
			expectedErr: "mismatch on array item index {0}. Details: field D is not matching. Source: 10, Target: 20",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := CompareItems(tc.source, tc.target, tc.depth)

			if !tc.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}
