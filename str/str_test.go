package str

import (
	"reflect"
	"testing"
)

func TestUint16Bytes(t *testing.T) {
	type args struct {
		u uint16
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "0",
			args: args{
				u: 0,
			},
			want: []byte{0, 0},
		},
		{
			name: "56",
			args: args{
				u: 56,
			},
			want: []byte{0, 56},
		},
		{
			name: "255",
			args: args{
				u: 255,
			},
			want: []byte{0, 255},
		},
		{
			name: "256",
			args: args{
				u: 256,
			},
			want: []byte{1, 0},
		},
		{
			name: "65535",
			args: args{
				u: 65535,
			},
			want: []byte{255, 255},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Uint16Bytes(tt.args.u); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint16Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesUint16(t *testing.T) {
	type args struct {
		bs []byte
	}
	tests := []struct {
		name string
		args args
		want uint16
	}{
		{
			name: "65535",
			args: args{
				bs: []byte{255, 255},
			},
			want: 65535,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesUint16(tt.args.bs); got != tt.want {
				t.Errorf("BytesUint16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNeat(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "",
		args: args{struct {
			Test string
			A    string
			B    struct{ C string }
		}{
			Test: "test",
			A:    "aa",
			B: struct{ C string }{
				C: "cc",
			},
		},
		},
		want: `{
  "Test": "test",
  "A": "aa",
  "B": {
    "C": "cc"
  }
}`,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Neat(tt.args.v); got != tt.want {
				t.Errorf("Neat() = %v, want %v", got, tt.want)
			}
		})
	}
}
