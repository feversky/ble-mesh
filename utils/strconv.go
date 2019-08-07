package utils

import "strconv"

func HexStringToUint(str string) uint {
	r, _ := strconv.ParseUint(str, 16, 32)
	return uint(r)
}

func UintToHexString(src uint) string {
	return strconv.FormatUint(uint64(src), 16)
}
