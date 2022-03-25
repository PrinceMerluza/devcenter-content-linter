package linter

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRefExistsCondition_Validate(t *testing.T) {
	tests := []struct {
		name      string
		condition *RefExistsCondition
		want      *ConditionResult
	}{
		{
			name: "Empty Markdown",
			condition: &RefExistsCondition{
				Path:              emptyFile,
				ReferencePatterns: &[]string{"!\\[.*\\]\\((.*)\\).*"},
			},
			want: &ConditionResult{
				IsSuccess:      true,
				FileHighlights: &[]FileHighlight{},
			},
		},
		{
			name: "Image Exists",
			condition: &RefExistsCondition{
				Path:              refExists,
				ReferencePatterns: &[]string{"!\\[.*\\]\\((.*)\\).*"},
			},
			want: &ConditionResult{
				IsSuccess: true,
				FileHighlights: &[]FileHighlight{
					{
						Path:        relPath(refExists),
						LineNumber:  6,
						LineCount:   1,
						LineContent: "![Image](yuri.png)",
					},
				},
			},
		},
		{
			name: "Image does not exist",
			condition: &RefExistsCondition{
				Path:              refExists2,
				ReferencePatterns: &[]string{"!\\[.*\\]\\((.*)\\).*"},
			},
			want: &ConditionResult{
				IsSuccess: false,
				FileHighlights: &[]FileHighlight{
					{
						Path:        relPath(refExists2),
						LineNumber:  6,
						LineCount:   1,
						LineContent: "![Image](yuri2.png)",
					},
				},
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

func TestRefExistsCondition_ValidateWithErrors(t *testing.T) {
	tests := []struct {
		name      string
		condition *RefExistsCondition
	}{
		{
			name: "Non-existent Path",
			condition: &RefExistsCondition{
				Path:              incorrectPath,
				ReferencePatterns: &[]string{"!\\[.*\\]\\((.*)\\).*"},
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
