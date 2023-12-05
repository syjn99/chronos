package helpers

import "github.com/prysmaticlabs/prysm/v4/math"

func BailOutRecoveryScore(valnum int) uint64 {
	num := uint64(valnum)
	switch {
	case num < math.PowerOf2(13): // 2^13
		return 396670510735312
	case num < math.PowerOf2(14): // 2^14
		return 424422339421802
	case num < math.PowerOf2(15): // 2^15
		return 445262211294194
	case num < math.PowerOf2(16): // 2^16
		return 460632171078136
	case num < math.PowerOf2(17): // 2^17
		return 471826519177477
	case num < math.PowerOf2(18): // 2^18
		return 479908445837684
	case num < math.PowerOf2(19): // 2^19
		return 485707547285690
	case num < math.PowerOf2(20): // 2^20
		return 489850697296405
	case num < math.PowerOf2(21): // 2^21
		return 492801774069274
	case num < math.PowerOf2(22): // 2^22
		return 494899265141145
	default:
		return 496387815679376
	}
}
