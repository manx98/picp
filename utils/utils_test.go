package utils

import "testing"

func TestByteSize(t *testing.T) {
	type args struct {
		bytes int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1.00KB",
			args: args{
				bytes: 1024,
			},
			want: "1.00KB",
		},
		{
			name: "1.00B",
			args: args{
				bytes: 1,
			},
			want: "1.00B",
		},
		{
			name: "0.98KB",
			args: args{
				bytes: 999,
			},
			want: "0.98KB",
		},
		{
			name: "1.95KB",
			args: args{
				bytes: 1999,
			},
			want: "1.95KB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ByteSize(tt.args.bytes); got != tt.want {
				t.Errorf("ByteSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
