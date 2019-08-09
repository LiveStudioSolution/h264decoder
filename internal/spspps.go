package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/32bitkid/bitreader"
	"github.com/LiveStudioSolution/h264decoder/internal/cavlc"
)

var errNotImplemented = fmt.Errorf("not implemented")

// h264 sequence parameters set
// T-REC-H.264-200711-S!!PDF-E.pdf 7.3.2.1.1 Sequence parameter set data syntax
type SPS struct {
	ProfileIdc         uint8
	ConstraintSet0Flag bool
	ConstraintSet1Flag bool
	ConstraintSet2Flag bool
	ConstraintSet3Flag bool
	ConstraintSet4Flag bool
	ConstraintSet5Flag bool

	LevelIdc uint8

	Id uint

	ChromaFormatIdc                 uint
	SeparateColourPlaneFlag         bool
	BitDepthLumaMinus8              uint
	BitDepthChromaMinus8            uint
	QpprimeYZeroTransformBypassFlag bool
	SeqScalingMatrixPresentFlag     bool
	SeqScalingListPresentFlag       []bool

	Log2MaxFrameNumMinus4          uint
	PicOrderCntType                uint
	Log2MaxPicOrderCntLsbMinus4L   uint
	DeltaPicOrderAlwaysZeroFlag    bool
	OffsetForNonRefPic             int
	OffsetForTopToBottomField      int
	NumRefFramesInPicOrderCntCycle uint
	OffsetForRefFrame              []int
	NumRefFrames                   uint

	GapsInFrameNumValueAllowedFlag bool
	PicWidthInMbsMinus1            uint
	PicHeightInMapUnitsMinus1      uint
	FrameMbsOnlyFlag               bool
	MbAdaptiveFrameFieldFlag       bool
	Direct8X8InferenceFlag         bool
	FrameCroppingFlag              bool
	FrameCrop                      FrameCrop
	VuiParametersPresentFlag       bool
	VuiParams                      VuiParameters

	trailingBits []byte

	br bitreader.BitReader
}

// T-REC-H.264-201402-S!!PDF-E.pdf Annex E.1.1 VUI parameters syntax
type VuiParameters struct {
	AspectRatioInfoPresentFlag bool
	AspectRatioIdc             uint8
	SarWidth                   uint16
	SarHeight                  uint16

	OverscanInfoPresentFlag      bool
	OverscanAppropriateFlag      bool
	VideoSignalTypePresentFlag   bool
	VideoFormat                  uint8
	VideoFullRangeFlag           bool
	ColourDescriptionPresentFlag bool
	ColourPrimaries              uint8
	TransferCharacteristics      uint8
	MatrixCoefficients           uint8

	ChromaLocInfoPresentFlag       bool
	ChromaSampleLocTypeTopField    uint
	ChromaSampleLocTypeBottomField uint

	TimingInfoPresentFlag bool
	NumUnitsInTick        uint32
	TimeScale             uint32
	FixedFrameRateFlag    bool

	NalHrdParametersPresentFlag bool
	VclHrdParametersPresentFlag bool
	HrdParameters               HrdParameters
	LowDelayHrdFlag             bool

	PicStructPresentFlag               bool
	BitstreamRestrictionFlag           bool
	MotionVectorsOverPicBoundariesFlag bool
	MaxBytesPerPicDenom                uint
	MaxBitsPerMbDenom                  uint
	Log2MaxMvLengthHorizontal          uint
	Log2MaxMvLengthVertical            uint
	MaxNumReorderFrames                uint
	MaxDecFrameBuffering               uint
}

type HrdParameters struct {
	CpbCntMinus1                       uint
	BitRateScale                       uint8
	CpbSizeScale                       uint8
	BitRateValueMinus1                 []uint
	CpbSizeValueMinus1                 []uint
	CbrFlag                            []bool
	InitialCpbRemovalDelayLengthMinus1 uint8
	CpbRemovalDelayLengthMinus1        uint8
	DpbOutputDelayLengthMinus1         uint8
	TimeOffsetLengths                  uint8
}

type FrameCrop struct {
	LeftOffset   uint
	RightOffset  uint
	TopOffset    uint
	BottomOffset uint
}

// h264 picture parameters set
type PPS struct {
}

func SpsParseFromRBSP(rbsp []byte) (*SPS, error) {
	sps := &SPS{}
	if err := sps.Load(rbsp); err != nil {
		return nil, err
	}
	return sps, nil
}

// ITU-T Recommendation H.264(200305)
//
func (sps *SPS) Load(rbsp []byte) error {
	sps.br = bitreader.NewReader(bytes.NewReader(rbsp))

	br := sps.br
	var err error
	sps.ProfileIdc, err = br.Read8(8)
	if err != nil {
		return err
	}
	sps.ConstraintSet0Flag, err = br.Read1()
	if err != nil {
		return err
	}
	sps.ConstraintSet1Flag, err = br.Read1()
	if err != nil {
		return err
	}
	sps.ConstraintSet2Flag, err = br.Read1()
	if err != nil {
		return err
	}
	sps.ConstraintSet3Flag, err = br.Read1()
	if err != nil {
		return err
	}
	sps.ConstraintSet4Flag, err = br.Read1()
	if err != nil {
		return err
	}
	sps.ConstraintSet5Flag, err = br.Read1()
	if err != nil {
		return err
	}

	// skip reserved 5 bits
	if err = br.Skip(2); err != nil {
		return err
	}

	sps.LevelIdc, err = br.Read8(8)
	if err != nil {
		return err
	}

	sps.Id, err = cavlc.DecUe(br)
	if err != nil {
		return err
	}

	sps.ChromaFormatIdc = 1 // default value 1

	if sps.ProfileIdc == 100 || sps.ProfileIdc == 110 || sps.ProfileIdc == 122 || sps.ProfileIdc == 244 ||
		sps.ProfileIdc == 44 || sps.ProfileIdc == 83 || sps.ProfileIdc == 86 || sps.ProfileIdc == 118 ||
		sps.ProfileIdc == 128 || sps.ProfileIdc == 138 || sps.ProfileIdc == 139 || sps.ProfileIdc == 134 {

		if err = sps.parseChromaFormat(); err != nil {
			return err
		}
	}

	sps.Log2MaxFrameNumMinus4, err = cavlc.DecUe(br)
	if err != nil {
		return err
	}

	sps.PicOrderCntType, err = cavlc.DecUe(br)
	if err != nil {
		return err
	}

	if sps.PicOrderCntType == 0 {
		sps.Log2MaxPicOrderCntLsbMinus4L, err = cavlc.DecUe(br)
		if err != nil {
			return err
		}
	} else if sps.PicOrderCntType == 1 {
		sps.DeltaPicOrderAlwaysZeroFlag, err = br.Read1()
		if err != nil {
			return err
		}
		sps.OffsetForNonRefPic, err = cavlc.DecSe(br)
		if err != nil {
			return err
		}
		sps.OffsetForTopToBottomField, err = cavlc.DecSe(br)
		if err != nil {
			return err
		}
		sps.NumRefFramesInPicOrderCntCycle, err = cavlc.DecUe(br)
		if err != nil {
			return err
		}
		if sps.NumRefFramesInPicOrderCntCycle > 0 {
			sps.OffsetForRefFrame = make([]int, sps.NumRefFramesInPicOrderCntCycle)
		}
		for i := uint(0); i < sps.NumRefFramesInPicOrderCntCycle; i++ {
			sps.OffsetForRefFrame[i], err = cavlc.DecSe(br)
		}

	}

	sps.NumRefFrames, err = cavlc.DecUe(br)
	if err != nil {
		return err
	}

	sps.GapsInFrameNumValueAllowedFlag, err = br.Read1()
	if err != nil {
		return err
	}

	sps.PicWidthInMbsMinus1, err = cavlc.DecUe(br)
	if err != nil {
		return err
	}
	sps.PicHeightInMapUnitsMinus1, err = cavlc.DecUe(br)
	if err != nil {
		return err
	}

	sps.FrameMbsOnlyFlag, err = br.Read1()
	if err != nil {
		return err
	}

	if !sps.FrameMbsOnlyFlag {
		sps.MbAdaptiveFrameFieldFlag, err = br.Read1()
		if err != nil {
			return err
		}
	}
	sps.Direct8X8InferenceFlag, err = br.Read1()
	if err != nil {
		return err
	}

	sps.FrameCroppingFlag, err = br.Read1()
	if err != nil {
		return err
	}

	if sps.FrameCroppingFlag {
		sps.FrameCrop.LeftOffset, err = cavlc.DecUe(br)
		if err != nil {
			return err
		}
		sps.FrameCrop.RightOffset, err = cavlc.DecUe(br)
		if err != nil {
			return err
		}
		sps.FrameCrop.TopOffset, err = cavlc.DecUe(br)
		if err != nil {
			return err
		}
		sps.FrameCrop.BottomOffset, err = cavlc.DecUe(br)
		if err != nil {
			return err
		}
	}

	sps.VuiParametersPresentFlag, err = br.Read1()
	if err != nil {
		return err
	}

	if sps.VuiParametersPresentFlag {
		return sps.parsingVuiParams()
	}

	return nil
}

func (sps *SPS) String() string {
	s, _ := json.Marshal(sps)
	return string(s)
}

const (
	ExtendedSAR = 255
)

// ITU-T Recommendation H.264(200305) Annex E
func (sps *SPS) parsingVuiParams() error {
	vui := &sps.VuiParams
	br := sps.br
	var err error

	if vui.AspectRatioInfoPresentFlag, err = br.Read1(); err != nil {
		return err
	}
	if vui.AspectRatioInfoPresentFlag {
		if vui.AspectRatioIdc, err = br.Read8(8); err != nil {
			return err
		}
		if vui.AspectRatioIdc == ExtendedSAR {
			if vui.SarWidth, err = br.Read16(16); err != nil {
				return err
			}
			if vui.SarHeight, err = br.Read16(16); err != nil {
				return err
			}
		}

	}

	if vui.OverscanInfoPresentFlag, err = br.Read1(); err != nil {
		return err
	}
	if vui.OverscanInfoPresentFlag {
		if vui.OverscanAppropriateFlag, err = br.Read1(); err != nil {
			return err
		}
	}

	if vui.VideoSignalTypePresentFlag, err = br.Read1(); err != nil {
		return err
	}
	if vui.VideoSignalTypePresentFlag {
		if vui.VideoFormat, err = br.Read8(3); err != nil {
			return err
		}
		if vui.VideoFullRangeFlag, err = br.Read1(); err != nil {
			return err
		}
		if vui.ColourDescriptionPresentFlag, err = br.Read1(); err != nil {
			return err
		}
		if vui.ColourDescriptionPresentFlag {
			if vui.ColourPrimaries, err = br.Read8(8); err != nil {
				return err
			}
			if vui.TransferCharacteristics, err = br.Read8(8); err != nil {
				return err
			}
			if vui.MatrixCoefficients, err = br.Read8(8); err != nil {
				return err
			}
		}
	}
	if vui.ChromaLocInfoPresentFlag, err = br.Read1(); err != nil {
		return err
	}
	if vui.ChromaLocInfoPresentFlag {
		if vui.ChromaSampleLocTypeTopField, err = cavlc.DecUe(br); err != nil {
			return err
		}
		if vui.ChromaSampleLocTypeBottomField, err = cavlc.DecUe(br); err != nil {
			return err
		}
	}

	if vui.TimingInfoPresentFlag, err = br.Read1(); err != nil {
		return err
	}
	if vui.TimingInfoPresentFlag {
		if vui.NumUnitsInTick, err = br.Read32(32); err != nil {
			return err
		}
		if vui.TimeScale, err = br.Read32(32); err != nil {
			return err
		}
		if vui.FixedFrameRateFlag, err = br.Read1(); err != nil {
			return err
		}
	}

	if vui.NalHrdParametersPresentFlag, err = br.Read1(); err != nil {
		return err
	}
	if vui.NalHrdParametersPresentFlag {
		if err = sps.parseHdrParameters(); err != nil {
			return err
		}
	}
	if vui.VclHrdParametersPresentFlag, err = br.Read1(); err != nil {
		return err
	}
	if vui.VclHrdParametersPresentFlag {
		if err = sps.parseHdrParameters(); err != nil {
			return err
		}
	}
	if vui.NalHrdParametersPresentFlag || vui.VclHrdParametersPresentFlag {
		if vui.LowDelayHrdFlag, err = br.Read1(); err != nil {
			return err
		}
	}
	if vui.PicStructPresentFlag, err = br.Read1(); err != nil {
		return err
	}
	if vui.BitstreamRestrictionFlag, err = br.Read1(); err != nil {
		return err
	}
	if vui.BitstreamRestrictionFlag {
		if vui.MotionVectorsOverPicBoundariesFlag, err = br.Read1(); err != nil {
			return err
		}
		if vui.MaxBytesPerPicDenom, err = cavlc.DecUe(br); err != nil {
			return err
		}
		if vui.MaxBitsPerMbDenom, err = cavlc.DecUe(br); err != nil {
			return err
		}
		if vui.Log2MaxMvLengthHorizontal, err = cavlc.DecUe(br); err != nil {
			return err
		}
		if vui.Log2MaxMvLengthVertical, err = cavlc.DecUe(br); err != nil {
			return err
		}
		if vui.MaxNumReorderFrames, err = cavlc.DecUe(br); err != nil {
			return err
		}
		if vui.MaxDecFrameBuffering, err = cavlc.DecUe(br); err != nil {
			return err
		}
	}
	return nil
}

func (sps *SPS) parseChromaFormat() error {
	// todo
	panic(errNotImplemented)
	return nil
}

func (sps *SPS) parseHdrParameters() error {
	hrd := &sps.VuiParams.HrdParameters
	br := sps.br
	var err error
	if hrd.CpbCntMinus1, err = cavlc.DecUe(br); err != nil {
		return err
	}
	if hrd.BitRateScale, err = br.Read8(4); err != nil {
		return err
	}
	if hrd.CpbSizeScale, err = br.Read8(4); err != nil {
		return err
	}
	hrd.BitRateValueMinus1 = make([]uint, hrd.CpbCntMinus1+1)
	hrd.CpbSizeValueMinus1 = make([]uint, hrd.CpbCntMinus1+1)
	hrd.CbrFlag = make([]bool, hrd.CpbCntMinus1+1)
	for sidx := uint(0); sidx <= hrd.CpbCntMinus1; sidx++ {
		if hrd.BitRateValueMinus1[sidx], err = cavlc.DecUe(br); err != nil {
			return err
		}
		if hrd.CpbSizeValueMinus1[sidx], err = cavlc.DecUe(br); err != nil {
			return err
		}
		if hrd.CbrFlag[sidx], err = br.Read1(); err != nil {
			return err
		}
	}
	if hrd.InitialCpbRemovalDelayLengthMinus1, err = br.Read8(5); err != nil {
		return err
	}
	if hrd.CpbRemovalDelayLengthMinus1, err = br.Read8(5); err != nil {
		return err
	}
	if hrd.DpbOutputDelayLengthMinus1, err = br.Read8(5); err != nil {
		return err
	}
	if hrd.TimeOffsetLengths, err = br.Read8(5); err != nil {
		return err
	}
	return nil
}
