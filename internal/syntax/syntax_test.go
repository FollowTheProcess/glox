package syntax_test

import (
	"fmt"
	"testing"

	"github.com/FollowTheProcess/glox/internal/syntax"
	"github.com/FollowTheProcess/test"
)

func TestPositionString(t *testing.T) {
	tests := []struct {
		name string          // Name of the test case
		want string          // Expected return value
		pos  syntax.Position // Position under test
	}{
		{
			name: "empty",
			pos:  syntax.Position{},
			want: `BadPosition: {Name: "", Line: 0, StartCol: 0, EndCol: 0}`,
		},
		{
			name: "missing name",
			pos:  syntax.Position{Line: 12, StartCol: 2, EndCol: 6},
			want: `BadPosition: {Name: "", Line: 12, StartCol: 2, EndCol: 6}`,
		},
		{
			name: "zero line",
			pos:  syntax.Position{Name: "file.txt", Line: 0, StartCol: 12, EndCol: 19},
			want: `BadPosition: {Name: "file.txt", Line: 0, StartCol: 12, EndCol: 19}`,
		},
		{
			name: "zero start column",
			pos:  syntax.Position{Name: "file.txt", Line: 4, StartCol: 0, EndCol: 19},
			want: `BadPosition: {Name: "file.txt", Line: 4, StartCol: 0, EndCol: 19}`,
		},
		{
			name: "zero end column",
			pos:  syntax.Position{Name: "file.txt", Line: 4, StartCol: 1, EndCol: 0},
			want: `BadPosition: {Name: "file.txt", Line: 4, StartCol: 1, EndCol: 0}`,
		},
		{
			name: "end less than start",
			pos:  syntax.Position{Name: "test.lox", Line: 1, StartCol: 6, EndCol: 4},
			want: `BadPosition: {Name: "test.lox", Line: 1, StartCol: 6, EndCol: 4}`,
		},
		{
			name: "valid single column",
			pos:  syntax.Position{Name: "demo.lox", Line: 1, StartCol: 6, EndCol: 6},
			want: "demo.lox:1:6",
		},
		{
			name: "valid column range",
			pos:  syntax.Position{Name: "demo.lox", Line: 17, StartCol: 20, EndCol: 26},
			want: "demo.lox:17:20-26",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Equal(t, tt.pos.String(), tt.want)
		})
	}
}

func FuzzPosition(f *testing.F) {
	f.Add("", 0, 0, 0)
	f.Add("name.txt", 1, 1, 2)
	f.Add("valid.lox", 12, 17, 19)
	f.Add("invalid.lox", 0, -9, 9999)

	f.Fuzz(func(t *testing.T, name string, line, startCol, endCol int) {
		pos := syntax.Position{
			Name:     name,
			Line:     line,
			StartCol: startCol,
			EndCol:   endCol,
		}

		got := pos.String()

		// Property: If IsValid returns false, the string must be this format
		if !pos.IsValid() {
			want := fmt.Sprintf(
				"BadPosition: {Name: %q, Line: %d, StartCol: %d, EndCol: %d}",
				name,
				line,
				startCol,
				endCol,
			)
			test.Equal(t, got, want)
			return
		}

		// Property: If IsValid returned true, Line must be >= 1
		test.True(
			t,
			pos.Line >= 1,
			test.Context("IsValid() = true but pos.Line (%d) was not >= 1", pos.Line),
		)

		// Property: If IsValid returned true, StartCol must be >= 1
		test.True(
			t,
			pos.StartCol >= 1,
			test.Context("IsValid() = true but pos.StartCol (%d) was not >= 1", pos.StartCol),
		)

		// Property: If IsValid returned true, EndCol must be >= 1
		test.True(
			t,
			pos.EndCol >= 1,
			test.Context("IsValid() = true but pos.EndCol (%d) was not >= 1", pos.EndCol),
		)

		// Property: If IsValid returned true, EndCol must also be >= StartCol
		test.True(
			t,
			pos.EndCol >= pos.StartCol,
			test.Context(
				"IsValid() = true but pos.EndCol (%d) was not >= pos.StartCol (%d)",
				pos.EndCol,
				pos.StartCol,
			),
		)

		// Property: If StartCol == EndCol, no range must appear in the string
		if startCol == endCol {
			want := fmt.Sprintf("%s:%d:%d", name, line, startCol)
			test.Equal(t, got, want)
			return
		}

		// Otherwise the position must be a valid position with a column range
		want := fmt.Sprintf("%s:%d:%d-%d", name, line, startCol, endCol)
		test.Equal(t, got, want)
	})
}
