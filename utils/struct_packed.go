package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
)

//https://gist.github.com/cloveryume/9a59e8d77f5836f11720#file-golang_struct_packed-go

// struct filed only support uint* and string
// for string, there should be a value to represent the length of string

// packed
func WriteStructToBuffer(data interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	v := reflect.Indirect(reflect.ValueOf(data))
	if v.Kind() != reflect.Struct {
		return nil, errors.New("invaild type Not a struct")
	}

	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Type().Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			err := binary.Write(buffer, binary.BigEndian, v.Field(i).Interface())
			if err != nil {
				return nil, err
			}
		case reflect.String:
			s := v.Field(i).String()
			b := []byte(s)
			_, err := buffer.Write(b)
			if err != nil {
				return nil, err
			}
		case reflect.Slice:
			s := v.Field(i).Bytes()
			err := binary.Write(buffer, binary.BigEndian, s)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("invail type Unspport reflect type")
		}
	}
	return buffer.Bytes(), nil
}

//unpacked
func ReadStructFromBuffer(b []byte, data interface{}) error {
	buffer := bytes.NewBuffer(b)
	k := reflect.TypeOf(data).Kind()
	if k != reflect.Ptr {
		return errors.New("the second parameter must be a pointer")
	}

	v := reflect.ValueOf(data).Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("invaild type Not a struct")
	}

	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Type().Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			err := binary.Read(buffer, binary.BigEndian, v.Field(i).Addr().Interface())
			if err != nil {
				return err
			}
		case reflect.String:
			var strlen uint16
			err := binary.Read(buffer, binary.BigEndian, &strlen)
			if err != nil {
				return err
			}

			v.Field(i).SetString(string(buffer.Next(int(strlen))))

		default:
			return errors.New("invail type Unspport reflect type")
		}
	}

	return nil
}
