package utils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDeleteAt(t *testing.T) {
	type args[T any] struct {
		src   []T
		index int
	}
	type testCase[T any] struct {
		name    string
		args    args[T]
		want    []T
		wantErr bool
	}
	tests := []testCase[any]{
		{
			name: "",
			args: args[any]{
				src:   []any{1, 2, 3, 4, 5},
				index: 3,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteAt(tt.args.src, tt.args.index)
			fmt.Println(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteAt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
