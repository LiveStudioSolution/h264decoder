package cavlc

import (
	"github.com/32bitkid/bitreader"
)

type CodedBlock struct {
	Intra44 uint32
	Inter   uint32
}

var codedBlockPatternMap []CodedBlock = []CodedBlock{
	// 0 ~ 4
	{47, 0},
	{31, 16},
	{15, 1},
	{0, 2},
	{23, 4},

	// 5~9
	{27, 8},
	{29, 32},
	{30, 3},
	{7, 5},
	{11, 10},

	// 10 ~ 14
	{13, 12},
	{14, 15},
	{39, 47},
	{43, 7},
	{45, 11},


	//15 ~ 19
	{46, 13},
	{16, 14},
	{3, 6},
	{5, 9},
	{10, 31},

	// 20 ~ 24
	{12, 35},
	{19, 37},
	{21, 42},
	{26, 44},
	{28, 33},

	// 25 ~ 29
	{35, 34},
	{37, 36},
	{42, 40},
	{44, 39},
	{1, 43},

	// 30 ~ 34
	{2, 45},
	{4, 46},
	{8, 17},
	{17, 18},
	{18, 20},

	// 35 ~ 39
	{20, 24},
	{24, 19},
	{6, 21},
	{9, 26},
	{22, 28},

	// 40 ~ 44
	{25, 23},
	{32, 27},
	{33, 29},
	{34, 30},
	{36, 22},

	// 45 ~ 49
	{40, 25},
	{38, 38},
	{41, 41},
}

func init() {

}

func DecUe(br bitreader.BitReader) (uint, error) {
	return readCodeNum(br)
}

func readCodeNum(br bitreader.BitReader) (uint, error) {
	var err error
	leadingZeroBits := -1
	for b := false; !b; leadingZeroBits++ {
		b, err = br.Read1()
		if err != nil {
			return 0, err
		}
	}
	var suffix uint32
	if leadingZeroBits > 0 {
		suffix, err = br.Read32(uint(leadingZeroBits))
		if err != nil {
			return 0, err
		}
	}
	//codeNum = 2^leadingZeroBits – 1 + read_bits( leadingZeroBits )
	codeNum := 1<<uint(leadingZeroBits) - 1 + int(suffix)
	return uint(codeNum), nil
}

func DecSe(br bitreader.BitReader) (int, error) {
	codeNum, err := DecUe(br)
	if err != nil {
		return 0, err
	}
	// k : codeNum
	//(–1)k+1 Ceil( k÷2 )
	v := (codeNum + 1) >> 1
	s := 1
	if codeNum&1 == 0 {
		s = -1
	}
	return s * int(v), nil
}

func DecMe(br bitreader.BitReader) (CodedBlock, error) {
	uv, err := DecUe(br)
	if err != nil {
		return CodedBlock{0, 0}, err
	}
	return codedBlockPatternMap[uv], nil
}

//func DecTe(br bitreader.BitReader)(uint,error) {
//	var err error
//	leadingZeroBits := -1
//	for b := false; !b; leadingZeroBits++ {
//		b, err = br.Read1()
//		if err != nil {
//			return 0, err
//		}
//	}
//	var suffix uint32
//	if leadingZeroBits > 0 {
//		suffix, err = br.Read32(uint(leadingZeroBits))
//		if err != nil {
//			return 0, err
//		}
//	}
//	//codeNum = 2^leadingZeroBits – 1 + read_bits( leadingZeroBits )
//	codeNum := 1<<uint(leadingZeroBits) - 1 + int(suffix)
//	return uint(codeNum), nil
//}
