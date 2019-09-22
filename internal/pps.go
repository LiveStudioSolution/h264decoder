package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/32bitkid/bitreader"
	"github.com/LiveStudioSolution/h264decoder/internal/rbr"
)

// h264 picture parameters set
// T-REC-H.264-201402-S!!PDF-E.pdf 7.3.2.2 Picture parameter set RBSP syntax
type PPS struct {
	Id                                    uint `json:"id"`
	SeqParameterSetId                     uint `json:"seq_parameter_set_id"`
	EntropyCodingModeFlag                 bool `json:"entropy_coding_mode_flag"`
	BottomFieldPicOrderInFramePresentFlag bool `json:"bottom_field_pic_order_in_frame_present_flag"`

	NumSliceGroupsMinus1          uint   `json:"num_slice_groups_minus_1"`
	SliceGroupMapType             uint   `json:"slice_group_map_type"`
	RunLengthMinus1               []uint `json:"run_length_minus_1"`
	TopLeft                       []uint `json:"top_left"`
	BottomRight                   []uint `json:"bottom_right"`
	SliceGroupChangeDirectionFlag bool   `json:"slice_group_change_direction_flag"`
	SliceGroupChangeRateMinus1    uint   `json:"slice_group_change_rate_minus_1"`
	PicSizeInMapUnitsMinus1       uint   `json:"pic_size_in_map_units_minus_1"`
	SliceGroupId                  []uint `json:"slice_group_id"`

	NumRefIdxL0DefaultActiveMinus1     uint  `json:"num_ref_idx_l_0_default_active_minus_1"`
	NumRefIdxL1DefaultActiveMinus1     uint  `json:"num_ref_idx_l_1_default_active_minus_1"`
	WeightedPredFlag                   bool  `json:"weighted_pred_flag"`
	WeightedBipredIdc                  uint8 `json:"weighted_bipred_idc"`
	PicInitQpMinus26                   int   `json:"pic_init_qp_minus_26"`
	PicInitQsMinus26                   int   `json:"pic_init_qs_minus_26"`
	ChromaQpIndexOffset                int   `json:"chroma_qp_index_offset"`
	DeblockingFilterControlPresentFlag bool  `json:"deblocking_filter_control_present_flag"`
	ConstrainedIntraPredFlag           bool  `json:"constrained_intra_pred_flag"`
	RedundantPicCntPresentFlag         bool  `json:"redundant_pic_cnt_present_flag"`

	Transform8X8ModeFlag        bool   `json:"transform_8_x_8_mode_flag"`
	PicScalingMatrixPresentFlag bool   `json:"pic_scaling_matrix_present_flag"`
	PicScalingListPresentFlag   []bool `json:"pic_scaling_list_present_flag"`
	SecondChromaQpIndexOffset   int    `json:"second_chroma_qp_index_offset"`

	br bitreader.BitReader
}

func (pps *PPS) Load(rbsp []byte) error {
	pps.br = bitreader.NewReader(bytes.NewReader(rbsp))
	br := pps.br
	var err error
	if pps.Id, err = rbr.DecUe(br); err != nil {
		return err
	}
	if pps.SeqParameterSetId, err = rbr.DecUe(br); err != nil {
		return err
	}
	if pps.EntropyCodingModeFlag, err = br.Read1(); err != nil {
		return err
	}

	if pps.BottomFieldPicOrderInFramePresentFlag, err = br.Read1(); err != nil {
		return err
	}

	if pps.NumSliceGroupsMinus1, err = rbr.DecUe(br); err != nil {
		return err
	}

	if pps.NumSliceGroupsMinus1 > 0 {
		if pps.SliceGroupMapType, err = rbr.DecUe(br); err != nil {
			return err
		}
		switch pps.SliceGroupMapType {
		case 0:
			pps.RunLengthMinus1 = make([]uint, pps.NumSliceGroupsMinus1+1)
			for iGroup := uint(0); iGroup <= pps.NumSliceGroupsMinus1; iGroup++ {
				if pps.RunLengthMinus1[iGroup], err = rbr.DecUe(br); err != nil {
					return err
				}
			}
		case 2:
			pps.TopLeft = make([]uint, pps.NumSliceGroupsMinus1+1)
			pps.BottomRight = make([]uint, pps.NumSliceGroupsMinus1+1)
			for iGroup := uint(0); iGroup <= pps.NumSliceGroupsMinus1; iGroup++ {
				if pps.TopLeft[iGroup], err = rbr.DecUe(br); err != nil {
					return err
				}
				if pps.BottomRight[iGroup], err = rbr.DecUe(br); err != nil {
					return err
				}
			}
		case 3, 4, 5:
			if pps.SliceGroupChangeDirectionFlag, err = br.Read1(); err != nil {
				return err
			}
			if pps.SliceGroupChangeRateMinus1, err = rbr.DecUe(br); err != nil {
				return err
			}
		case 6:
			if pps.PicSizeInMapUnitsMinus1, err = rbr.DecUe(br); err != nil {
				return err
			}
			pps.SliceGroupId = make([]uint, pps.PicSizeInMapUnitsMinus1+1)
			for iGroup := uint(0); iGroup <= pps.PicSizeInMapUnitsMinus1; iGroup++ {
				if pps.SliceGroupId[iGroup], err = rbr.DecUe(br); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("error SliceGroupMapType %v", pps.SliceGroupMapType)
		}
	}

	if pps.NumRefIdxL0DefaultActiveMinus1, err = rbr.DecUe(br); err != nil {
		return err
	}
	if pps.NumRefIdxL1DefaultActiveMinus1, err = rbr.DecUe(br); err != nil {
		return err
	}
	if pps.WeightedPredFlag, err = br.Read1(); err != nil {
		return err
	}
	if pps.WeightedBipredIdc, err = br.Read8(2); err != nil {
		return err
	}
	if pps.PicInitQpMinus26, err = rbr.DecSe(br); err != nil {
		return err
	}
	if pps.PicInitQsMinus26, err = rbr.DecSe(br); err != nil {
		return err
	}
	if pps.ChromaQpIndexOffset, err = rbr.DecSe(br); err != nil {
		return err
	}
	if pps.DeblockingFilterControlPresentFlag, err = br.Read1(); err != nil {
		return err
	}
	if pps.ConstrainedIntraPredFlag, err = br.Read1(); err != nil {
		return err
	}
	if pps.RedundantPicCntPresentFlag, err = br.Read1(); err != nil {
		return err
	}
	// todo MoreRBSPDatao
	//if ! rbr.MoreRBSPData(br) {
	//	return nil
	//}
	return nil
}

func ParsePpsFromRBSP(rbsp []byte) (*PPS, error) {
	pps := &PPS{}
	if err := pps.Load(rbsp); err != nil {
		return nil, err
	}
	return pps, nil
}

func (pps *PPS) String() string {
	s, _ := json.Marshal(pps)
	return string(s)
}
