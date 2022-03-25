package linter

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestNotContainsCondition_Validate(t *testing.T) {
	tests := []struct {
		name      string
		condition *NotContainsCondition
		want      *ConditionResult
	}{
		{
			name: "Empty Markdown",
			condition: &NotContainsCondition{
				Path:        emptyFile,
				NotContains: &[]string{"RANDOM"},
			},
			want: &ConditionResult{
				IsSuccess:      true,
				FileHighlights: &[]FileHighlight{},
			},
		},
		{
			name: "Link with no alternate text",
			condition: &NotContainsCondition{
				Path:        notContainsFile,
				NotContains: &[]string{"\\[.*\\]\\(.*[^ \"]+[^\"]*\\)"},
			},
			want: &ConditionResult{
				IsSuccess: false,
				FileHighlights: &[]FileHighlight{
					{
						Path:        relPath(notContainsFile),
						LineNumber:  5,
						LineCount:   1,
						LineContent: "bbbbb [Link Example](b/b)",
					},
				},
			},
		},
		{
			name: "Random Unique Text",
			condition: &NotContainsCondition{
				Path:        notContainsFile,
				NotContains: &[]string{uuid.NewString()},
			},
			want: &ConditionResult{
				IsSuccess:      true,
				FileHighlights: &[]FileHighlight{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.condition.Validate(); !cmp.Equal(got, tt.want) {
				t.Errorf("%v", cmp.Diff(got, tt.want))
				if got.Error != nil {
					t.Errorf("Error: %v", got.Error)
				}
			}
		})
	}
}

func TestNotContainsCondition_ValidateWithErrors(t *testing.T) {
	tests := []struct {
		name      string
		condition *NotContainsCondition
		want      *ConditionResult
	}{
		{
			name: "Empty regex",
			condition: &NotContainsCondition{
				Path:        notContainsFile,
				NotContains: &[]string{""},
			},
		},
		{
			name: "Incorrect Path",
			condition: &NotContainsCondition{
				Path:        incorrectPath,
				NotContains: &[]string{"RANDOM"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.condition.Validate(); got.Error == nil {
				t.Errorf("Expected error, got nil")
			}
		})
	}
}
