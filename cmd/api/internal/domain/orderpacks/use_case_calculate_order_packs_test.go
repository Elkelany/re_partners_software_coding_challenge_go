package orderpacks

import (
	"reflect"
	"testing"
)

func TestUseCaseCalculateOrderPacks_Run(t *testing.T) {
	type args struct {
		request UseCaseCalculateOrderPacksRequest
	}

	tests := []struct {
		name    string
		args    args
		want    map[uint64]uint64
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				request: UseCaseCalculateOrderPacksRequest{
					PackSizes:  []uint64{250, 500, 1000, 2000, 5000},
					OrderItems: 1,
				},
			},
			want:    map[uint64]uint64{250: 1, 500: 0, 1000: 0, 2000: 0, 5000: 0},
			wantErr: false,
		},
		{
			name: "99",
			args: args{
				request: UseCaseCalculateOrderPacksRequest{
					PackSizes:  []uint64{11, 15},
					OrderItems: 99,
				},
			},
			want:    map[uint64]uint64{11: 9, 15: 0},
			wantErr: false,
		},
		{
			name: "100",
			args: args{
				request: UseCaseCalculateOrderPacksRequest{
					PackSizes:  []uint64{11, 15, 17},
					OrderItems: 100,
				},
			},
			want:    map[uint64]uint64{15: 1, 17: 5},
			wantErr: false,
		},
		{
			name: "250",
			args: args{
				request: UseCaseCalculateOrderPacksRequest{
					PackSizes:  []uint64{250, 500, 1000, 2000, 5000},
					OrderItems: 250,
				},
			},
			want:    map[uint64]uint64{250: 1, 500: 0, 1000: 0, 2000: 0, 5000: 0},
			wantErr: false,
		},
		{
			name: "251",
			args: args{
				request: UseCaseCalculateOrderPacksRequest{
					PackSizes:  []uint64{250, 500, 1000, 2000, 5000},
					OrderItems: 251,
				},
			},
			want:    map[uint64]uint64{250: 0, 500: 1, 1000: 0, 2000: 0, 5000: 0},
			wantErr: false,
		},
		{
			name: "501",
			args: args{
				request: UseCaseCalculateOrderPacksRequest{
					PackSizes:  []uint64{250, 500, 1000, 2000, 5000},
					OrderItems: 501,
				},
			},
			want:    map[uint64]uint64{250: 1, 500: 1, 1000: 0, 2000: 0, 5000: 0},
			wantErr: false,
		},
		{
			name: "599",
			args: args{
				request: UseCaseCalculateOrderPacksRequest{
					PackSizes:  []uint64{200, 600, 1000, 2000, 5000},
					OrderItems: 599,
				},
			},
			want:    map[uint64]uint64{200: 0, 600: 1, 1000: 0, 2000: 0, 5000: 0},
			wantErr: false,
		},
		{
			name: "12001",
			args: args{
				request: UseCaseCalculateOrderPacksRequest{
					PackSizes:  []uint64{250, 500, 1000, 2000, 5000},
					OrderItems: 12001,
				},
			},
			want:    map[uint64]uint64{250: 1, 500: 0, 1000: 0, 2000: 1, 5000: 2},
			wantErr: false,
		},
		{
			name: "12251",
			args: args{
				request: UseCaseCalculateOrderPacksRequest{
					PackSizes:  []uint64{250, 500, 1000, 2000, 5000},
					OrderItems: 12251,
				},
			},
			want:    map[uint64]uint64{250: 0, 500: 1, 1000: 0, 2000: 1, 5000: 2},
			wantErr: false,
		},
		{
			name: "500000",
			args: args{
				request: UseCaseCalculateOrderPacksRequest{
					PackSizes:  []uint64{23, 31, 53},
					OrderItems: 500000,
				},
			},
			want:    map[uint64]uint64{23: 2, 31: 7, 53: 9429},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				request: UseCaseCalculateOrderPacksRequest{
					PackSizes:  []uint64{1, 250, 0, 500},
					OrderItems: 251,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UseCaseCalculateOrderPacks{}

			got, err := u.Run(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() got = %v, want %v", got, tt.want)
			}
		})
	}
}
