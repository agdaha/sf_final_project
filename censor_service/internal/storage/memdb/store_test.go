package memdb

import (
	"reflect"
	"testing"
)

func TestStore_Words(t *testing.T) {
	type fields struct {
		db []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name:   "memdb",
			fields: fields{db: []string{"qwerty", "йцукен", "zxcvbn", "asdfgh"}},
			want:   []string{"qwerty", "йцукен", "zxcvbn", "asdfgh"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				db: tt.fields.db,
			}
			if got := s.Words(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Store.Words() = %v, want %v", got, tt.want)
			}
		})
	}
}
