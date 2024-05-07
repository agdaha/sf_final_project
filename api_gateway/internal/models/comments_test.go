package models

import (
	"database/sql"
	"reflect"
	"testing"
)

func TestNullInt64_MarshalJSON(t *testing.T) {
	type fields struct {
		NullInt64 sql.NullInt64
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name:    "marshalJSON valid",
			fields:  fields{sql.NullInt64{Int64: 0, Valid: true}},
			want:    []byte("0"),
			wantErr: false,
		},
		{
			name:    "marshalJSON null",
			fields:  fields{sql.NullInt64{Int64: 0, Valid: false}},
			want:    []byte("null"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ni := &NullInt64{
				NullInt64: tt.fields.NullInt64,
			}
			got, err := ni.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("NullInt64.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NullInt64.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullInt64_UnmarshalJSON(t *testing.T) {
	type fields struct {
		NullInt64 sql.NullInt64
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Unmarshal 0",
			fields:  fields{NullInt64: sql.NullInt64{Int64: 0, Valid: true}},
			args:    args{data: []byte("0")},
			wantErr: false,
		},
		{
			name:    "Unmarshal null",
			fields:  fields{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
			args:    args{data: []byte("null")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &NullInt64{
				NullInt64: tt.fields.NullInt64,
			}
			if err := v.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("NullInt64.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
