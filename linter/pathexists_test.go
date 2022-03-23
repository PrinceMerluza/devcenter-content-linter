package linter

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPathExistsCondition_Validate(t *testing.T) {
	tests := []struct {
		name      string
		condition *PathExistsCondition
		want      *ConditionResult
	}{
		{
			name: "Empty Path",
			condition: &PathExistsCondition{
				Path: "",
			},
			want: &ConditionResult{
				IsSuccess: false,
			},
		},
		{
			name: "Existent Path",
			condition: &PathExistsCondition{
				Path: emptyFile,
			},
			want: &ConditionResult{
				IsSuccess: true,
			},
		},
		{
			name: "Non-Existent Path",
			condition: &PathExistsCondition{
				Path: incorrectPath,
			},
			want: &ConditionResult{
				IsSuccess: false,
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
