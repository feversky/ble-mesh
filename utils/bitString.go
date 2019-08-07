package utils

import (
	"encoding/binary"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"ble-mesh/utils/errors"

	funk "github.com/thoas/go-funk"
)

const (
	// 5,13,08,3
	// 5,11,8/16
	// 5,11,B8
	unPackRegex = `^(\d+?\s*,\s*)*((\d+\s*)|(B[1-9]+)|((\d+/)+\d+))$`
	packRegex   = `^((B\d+,)|(\d+?\s*,\s*))*(B\d+|((\d+/)?\d+\s*))$`
)

func InitializeStruct(v reflect.Value, level int) {
	if level == 0 {
		return
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)
		switch ft.Type.Kind() {
		case reflect.Map:
			f.Set(reflect.MakeMap(ft.Type))
		case reflect.Slice:
			f.Set(reflect.MakeSlice(ft.Type, 0, 0))
		case reflect.Chan:
			f.Set(reflect.MakeChan(ft.Type, 0))
		case reflect.Struct:
			InitializeStruct(f, level-1)
		case reflect.Ptr:
			fv := reflect.New(ft.Type.Elem())
			InitializeStruct(fv.Elem(), level-1)
			f.Set(fv)
		default:
		}
	}
}

func UnpackStructBE(data []byte, out interface{}) error {
	return unpackStruct(true, data, out)
}

func UnpackStructLE(data []byte, out interface{}) error {
	return unpackStruct(false, data, out)
}

func parseStructBits(v reflect.Value, pack bool) ([]interface{}, string, error) {
	bitString := ""
	fields := make([]interface{}, 0)
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get("bits")
		kind := v.Type().Field(i).Type.Kind()
		if kind == reflect.Slice {
			tag = "B" + tag
		} else if kind == reflect.Struct {
			f, s, err := parseStructBits(v.Field(i), pack)
			if err != nil {
				return nil, "", err
			}
			fields = append(fields, f...)
			bitString = bitString + s
			continue
		}
		bitString += tag + ","
		if tag[0] == '0' {
			continue
		} else {
			if !pack {
				// return pointer to filed, todo: handle the exception
				if v.Field(i).CanAddr() {
					fields = append(fields, v.Field(i).Addr().Interface())
				} else {
					return nil, "", errors.FieldOfStructNotAddressable.New().AddContextF("%+#v", v.Field(i))
				}
			} else {
				fields = append(fields, v.Field(i).Interface())
			}
		}
	}
	return fields, bitString, nil
}

func unpackStruct(be bool, data []byte, out interface{}) error {
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		return errors.NotAPointer.New().AddContextF("output type must be pointer to struct " + reflect.TypeOf(out).Kind().String())
	}
	v := reflect.Indirect(reflect.ValueOf(out).Elem())
	InitializeStruct(v, 5)

	fields, bitString, err := parseStructBits(v, false)
	if err != nil {
		return err
	}
	// remove the final ','
	bitString = bitString[:len(bitString)-1]
	// check length of optional fields
	lenTotalBits := int64(0)
	if opt := strings.IndexRune(bitString, 't'); opt != -1 {
		fs := strings.Split(bitString[:opt-1], ",")
		for _, f := range fs {
			d, _ := strconv.ParseInt(f, 10, 16)
			lenTotalBits += d
		}
		if int(lenTotalBits)/8 == len(data) {
			bitString = bitString[:opt-1]
			fields = fields[:len(fs)]
		}
	}

	// fmt.Printf("%+#v", bitString)
	return unpack(be, data, bitString, fields...)
}

func PackStructBE(in interface{}) (out []byte, err error) {
	return packStruct(true, in)
}

func PackStructLE(in interface{}) (out []byte, err error) {
	return packStruct(false, in)
}

func packStruct(be bool, in interface{}) (out []byte, err error) {
	v := reflect.ValueOf(in)
	if v.Type().Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}

	fields, bitString, err := parseStructBits(v, true)
	if err != nil {
		return nil, err
	}
	bitString = bitString[:len(bitString)-1]
	// fmt.Println(fields, bitString)
	return pack(be, bitString, fields...)
}

func PackBE(bitString string, in ...interface{}) (out []byte, err error) {
	return pack(true, bitString, in...)
}

func PackLE(bitString string, in ...interface{}) (out []byte, err error) {
	return pack(false, bitString, in...)
}

func value2bytes(be bool, value uint64, startBit, lenBit uint) []byte {
	mask := uint64(0)
	if be {
		mask = value << uint((8-startBit%8-lenBit%8)%8)
	} else {
		mask = value << uint(startBit%8)
	}

	nBytes := int(math.Ceil(float64(startBit%8+lenBit) / 8))
	b := make([]byte, 8)
	if nBytes == 1 {
		b[0] = byte(mask)
		b = b[0:1]
	} else if nBytes == 2 {
		if be {
			binary.BigEndian.PutUint16(b, uint16(mask))
		} else {
			binary.LittleEndian.PutUint16(b, uint16(mask))
		}
		b = b[:2]
	} else if nBytes <= 4 {
		if be {
			binary.BigEndian.PutUint32(b, uint32(mask))
			b = b[4-nBytes:]
		} else {
			binary.LittleEndian.PutUint32(b, uint32(mask))
			b = b[:nBytes]
		}
	} else if nBytes <= 8 {
		if be {
			binary.BigEndian.PutUint64(b, uint64(mask))
			b = b[8-nBytes:]
		} else {
			binary.LittleEndian.PutUint64(b, uint64(mask))
			b = b[:nBytes]
		}
	}
	return b
}

func pack(be bool, bitString string, in ...interface{}) (out []byte, err error) {
	match, _ := regexp.MatchString(packRegex, bitString)
	if !match {
		return nil, errors.WrongFormatOfBitString.New().AddContextF(bitString)
	}
	strs := strings.Split(bitString, ",")

	// special handling of model identifier
	sep := "/"
	strLast := strs[len(strs)-1]
	if strings.Contains(strLast, sep) {
		val := reflect.ValueOf(in[len(in)-1])
		conv := val.Convert(reflect.TypeOf(uint64(0)))
		actual := conv.Uint()
		for _, opt := range strings.Split(strLast, sep) {
			size, _ := strconv.ParseInt(opt, 10, 16)
			if actual <= (2 << uint(size)) {
				strs[len(strs)-1] = opt
				break
			}
		}
	}
	index := 0
	lenTotalBits := int64(0)
	startBit := []uint{}
	lenBit := []uint{}
	values := []uint64{}
	for _, str := range strs {
		str = strings.TrimSpace(str)
		if str[0] == 'B' {
			sli := reflect.ValueOf(in[index])
			d, _ := strconv.ParseInt(str[1:], 10, 16)
			for i := 0; i < sli.Len(); i++ {
				v := sli.Index(i).Convert(reflect.TypeOf(uint64(0))).Uint()
				values = append(values, v)
				startBit = append(startBit, uint(lenTotalBits))
				lenTotalBits += d
				lenBit = append(lenBit, uint(d))
			}
			index++
		} else {
			d, _ := strconv.ParseInt(str, 10, 16)
			startBit = append(startBit, uint(lenTotalBits))
			lenTotalBits += d
			lenBit = append(lenBit, uint(d))
			var actVal uint64
			if str[0] != '0' {
				actVal = reflect.ValueOf(in[index]).Convert(reflect.TypeOf(uint64(0))).Uint()
				index++
			}
			values = append(values, actVal)
		}
	}

	// if lenTotalBits%8 != 0 {
	// 	return nil, errors.New("Bit length incorrect")
	// }
	buf := make([]byte, lenTotalBits/8, lenTotalBits/8)
	for index := 0; index < len(startBit); index++ {
		b := value2bytes(be, values[index], startBit[index], lenBit[index])

		nBytes := int(math.Ceil(float64(startBit[index]%8+lenBit[index]) / 8))
		startByte := startBit[index] / 8
		for j := 0; j < nBytes; j++ {
			buf[startByte+uint(j)] |= b[j]
		}
	}

	return buf, nil
}

//
func pack1(be bool, bitString string, in ...interface{}) (out []byte, err error) {
	match, _ := regexp.MatchString(packRegex, bitString)
	if !match {
		return nil, errors.WrongFormatOfBitString.New().AddContextF(bitString)
	}
	strs := strings.Split(bitString, ",")
	lenTotalBits := int64(0)
	index := 0
	startBit := make([]uint, len(in))
	lenBit := make([]uint, len(in))
	typ := make([]rune, len(in))
	for _, str := range strs {
		str = strings.TrimSpace(str)
		if str[0] == '0' {
			d, _ := strconv.ParseInt(str, 10, 16)
			lenTotalBits += d
		} else if str[0] == 'B' {
			sli := in[index].([]byte)
			startBit[index] = uint(lenTotalBits)
			bitLen64, _ := strconv.ParseInt(str[1:], 10, 16)
			// if lenTotalBits%8 != 0 {
			// 	return nil, errors.New("start bit of B is incorrect")
			// }
			lenBit[index] = uint(len(sli) * int(bitLen64))
			lenTotalBits += int64(lenBit[index])
			typ[index] = 'B'
			index++
		} else {
			d, _ := strconv.ParseInt(str, 10, 16)
			startBit[index] = uint(lenTotalBits)
			lenTotalBits += d
			lenBit[index] = uint(d)
			typ[index] = 'b'
			index++
		}
	}
	// if lenTotalBits%8 != 0 {
	// 	return nil, errors.New("Bit length incorrect")
	// }
	buf := make([]byte, lenTotalBits/8, lenTotalBits/8)
	for index := 0; index < len(startBit); index++ {
		if typ[index] == 'B' {
			sli := in[index].([]byte)
			copy(buf[startBit[index]/8:], sli)
		} else if typ[index] == 'b' {
			b := make([]byte, 8)
			val := reflect.ValueOf(in[index])
			if val.Type().Kind() == reflect.Ptr {
				val = reflect.Indirect(val)
			}
			conv := val.Convert(reflect.TypeOf(uint64(0)))
			actual := conv.Uint()
			mask := uint64(0)
			if be {
				mask = actual << uint((8-startBit[index]%8-lenBit[index]%8)%8)
			} else {
				mask = actual << uint(startBit[index]%8)
			}

			nBytes := int(math.Ceil(float64(startBit[index]%8+lenBit[index]) / 8))
			startByte := startBit[index] / 8
			if nBytes == 1 {
				b[0] = byte(mask)
				b = b[0:1]
			} else if nBytes == 2 {
				if be {
					binary.BigEndian.PutUint16(b, uint16(mask))
				} else {
					binary.LittleEndian.PutUint16(b, uint16(mask))
				}
				b = b[:2]
			} else if nBytes <= 4 {
				if be {
					binary.BigEndian.PutUint32(b, uint32(mask))
					b = b[4-nBytes:]
				} else {
					binary.LittleEndian.PutUint32(b, uint32(mask))
					b = b[:nBytes]
				}
			} else if nBytes <= 8 {
				if be {
					binary.BigEndian.PutUint64(b, uint64(mask))
					b = b[8-nBytes:]
				} else {
					binary.LittleEndian.PutUint64(b, uint64(mask))
					b = b[:nBytes]
				}
			}
			for j := 0; j < nBytes; j++ {
				buf[startByte+uint(j)] |= b[j]
			}
		}
	}

	return buf, nil
}

// Bytes are not supported yet
func UnpackLE(input []byte, bitString string, out ...interface{}) error {
	return unpack(false, input, bitString, out...)
}

// Bytes are not supported yet
func UnpackBE(input []byte, bitString string, out ...interface{}) error {
	return unpack(true, input, bitString, out...)
}

func unpack(be bool, input []byte, bitString string, out ...interface{}) error {
	match, _ := regexp.MatchString(unPackRegex, bitString)
	if !match {
		return errors.WrongFormatOfBitString.New().AddContextF(bitString)
	}
	strs := strings.Split(bitString, ",")
	lastSeg := strings.TrimSpace(strs[len(strs)-1])
	sum := 0
	for _, str := range strs[:len(strs)-1] {
		str = strings.TrimSpace(str)
		bitLen64, _ := strconv.ParseInt(str, 10, 16)
		sum += int(bitLen64)
	}
	nRestBits := len(input)*8 - sum

	sep := "/"
	vl := strings.Split(lastSeg, sep)
	if funk.ContainsString(vl, strconv.Itoa(nRestBits)) {
		strs[len(strs)-1] = strconv.Itoa(nRestBits)
	} else if lastSeg[0] == 'B' {
		bitLen64, _ := strconv.ParseInt(lastSeg[1:], 10, 16)
		lenSlice := nRestBits / int(bitLen64)
		strs = strs[:len(strs)-1]
		for i := 0; i < lenSlice; i++ {
			strs = append(strs, lastSeg[1:])
		}
	} else {
		return errors.LengthMismatchOfBitString.New().AddContextF("data: %+#v, bitsting: %s", input, bitString)
	}

	temp := []uint{}
	bitOffset := 0
	for _, str := range strs {
		str = strings.TrimSpace(str)
		if str[0] == '0' {
			padding, _ := strconv.ParseInt(str, 10, 16)
			bitOffset += int(padding)
			continue
		}

		bitLen64, _ := strconv.ParseInt(str, 10, 16)
		bitLen := int(bitLen64)
		if bitLen > 64 {
			return errors.BitLenOver64OfBitString.New().AddContextF(bitString)
		}
		byteOffsetStart := bitOffset / 8
		byteOffsetEnd := (bitOffset + bitLen) / 8
		if (bitOffset+bitLen)%8 > 0 {
			byteOffsetEnd++
		}
		bytes := input[byteOffsetStart:byteOffsetEnd]
		var value uint64
		for i := 0; i < len(bytes); i++ {
			if be {
				value |= uint64(int(bytes[i]) << (uint(len(bytes)-i-1) * 8))
			} else {
				value |= uint64(int(bytes[i]) << uint(i*8))
			}
		}
		if be {
			value = value >> uint(len(bytes)*8-bitOffset%8-bitLen)
		} else {
			value = value >> uint(bitOffset%8)
		}
		mask := uint64(1)<<uint(bitLen) - 1
		value = value & mask
		// if be {
		// 	value = value << uint(64-len(bytes)*8+bitOffset%8)
		// 	value = value >> uint(64-bitLen)
		// } else {
		// 	value = value << uint(64-len(bytes)*8+bitOffset%8+8-bitLen%8)
		// 	value = value >> uint(64-bitLen)
		// }

		temp = append(temp, uint(value))
		bitOffset += bitLen
	}
	if len(temp) < len(out) && lastSeg[0] != 'B' {
		return errors.OutLenMismatchOfBitString.New().AddContextF(bitString)
	}

	for i := 0; i < len(out)-1; i++ {
		outType := reflect.TypeOf(out[i])

		if outType.Kind() != reflect.Ptr {
			return errors.NotAPointer.New().AddContextF("out type: %s", outType)
		}
		refVal := reflect.ValueOf(out[i])
		refVal.Elem().Set(reflect.ValueOf(temp[i]).Convert(reflect.Indirect(refVal).Type()))
	}
	lastOut := out[len(out)-1]
	outType := reflect.TypeOf(lastOut).Kind()
	refVal := reflect.ValueOf(lastOut)
	if outType == reflect.Ptr {
		if refVal.Elem().Type().Kind() == reflect.Slice {
			for i := len(out) - 1; i < len(temp); i++ {
				typ := refVal.Elem().Type().Elem()
				refVal.Elem().Set(reflect.Append(refVal.Elem(), reflect.ValueOf(temp[i]).Convert(typ)))
			}
		} else {
			refVal.Elem().Set(reflect.ValueOf(temp[len(temp)-1]).Convert(reflect.Indirect(refVal).Type()))
		}
	} else {
		return errors.NotAPointer.New().AddContextF("out type: %s", outType)
	}
	return nil
}
