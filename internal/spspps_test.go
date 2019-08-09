package internal

import (
	"reflect"
	"testing"
)

func TestSpsParseFromRBSP(t *testing.T) {
	type args struct {
		rbsp []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *SPS
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"1",
			args{
				[]byte{0x42, 0xC0, 0x1E, 0xDA, 0x02, 0x80, 0xBF, 0xE5, 0xC0,
					0x5A, 0x80, 0x80, 0x80, 0xA0, 0x00, 0x00, 0x7D, 0x20, 0x00, 0x1D, 0x4C, 0x01, 0xE2, 0xC5},
			},
			&SPS{
				ProfileIdc:66,
				ConstraintSet0Flag:true,
				ConstraintSet1Flag:true,
				ConstraintSet2Flag:false,
				ConstraintSet3Flag:false,
				ConstraintSet4Flag:false,
				ConstraintSet5Flag:false,

				LevelIdc:30,
				Id:0,
				ChromaFormatIdc:1,
				// mis

				Log2MaxFrameNumMinus4:0,
				PicOrderCntType:2,
				Log2MaxPicOrderCntLsbMinus4L:0,
				DeltaPicOrderAlwaysZeroFlag:false,
				OffsetForNonRefPic:0,
				OffsetForTopToBottomField:0,
				NumRefFramesInPicOrderCntCycle:0,

				NumRefFrames:1,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SpsParseFromRBSP(tt.args.rbsp)
			if (err != nil) != tt.wantErr {
				t.Errorf("SpsParseFromRBSP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SpsParseFromRBSP() = %v, want %v", got, tt.want)
			}
		})
	}
}
