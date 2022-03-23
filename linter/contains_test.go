package linter

import (
	"testing"

	"github.com/PrinceMerluza/devcenter-content-linter/config"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestContainsCondition_Validate(t *testing.T) {
	tests := []struct {
		name      string
		condition *ContainsCondition
		want      *ConditionResult
	}{
		{
			name: "Empty Markdown",
			condition: &ContainsCondition{
				Path: emptyFile,
				ContainsArr: &[]config.ContainsCondition{
					{
						Type:  "static",
						Value: uuid.New().String(),
					},
				},
			},
			want: &ConditionResult{
				IsSuccess:      false,
				FileHighlights: &[]FileHighlight{},
			},
		},
		{
			name: "Finding Waldo",
			condition: &ContainsCondition{
				Path: containsFile,
				ContainsArr: &[]config.ContainsCondition{
					{
						Type:  "static",
						Value: "WALDO",
					},
				},
			},
			want: &ConditionResult{
				IsSuccess: true,
				FileHighlights: &[]FileHighlight{
					{
						Path:        containsFile,
						LineNumber:  3,
						LineCount:   1,
						LineContent: "Laboris ea elit voluptate WALDO ullamco esse in fugiat ullamco",
					},
				},
			},
		},
		{
			name: "Multiple: Waldo + Regex",
			condition: &ContainsCondition{
				Path: containsFile,
				ContainsArr: &[]config.ContainsCondition{
					{
						Type:  "static",
						Value: "WALDO",
					},
					{
						Type:  "regex",
						Value: "## something.*",
					},
				},
			},
			want: &ConditionResult{
				IsSuccess: true,
				FileHighlights: &[]FileHighlight{
					{
						Path:        containsFile,
						LineNumber:  3,
						LineCount:   1,
						LineContent: "Laboris ea elit voluptate WALDO ullamco esse in fugiat ullamco",
					},
					{
						Path:        containsFile,
						LineNumber:  7,
						LineCount:   1,
						LineContent: "## something random text random",
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

func TestContainsCondition_ValidateWithErrors(t *testing.T) {
	tests := []struct {
		name      string
		condition *ContainsCondition
		want      *ConditionResult
	}{
		{
			name: "Empty static or regex",
			condition: &ContainsCondition{
				Path: containsFile,
				ContainsArr: &[]config.ContainsCondition{
					{
						Type:  "static",
						Value: "",
					},
					{
						Type:  "regex",
						Value: "",
					},
				},
			},
		},
		{
			name: "Incorrect Path",
			condition: &ContainsCondition{
				Path: incorrectPath,
				ContainsArr: &[]config.ContainsCondition{
					{
						Type:  "static",
						Value: "WALDO",
					},
				},
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
