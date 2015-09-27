package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

// Unmarshaller allows you to customize the unmarshal process of a field in CSV
// file
type Unmarshaller interface {
	UnmarshalCSV(string) error
}

var (
	intKindToSize = map[reflect.Kind]int{
		reflect.Int:   0,
		reflect.Int8:  8,
		reflect.Int16: 16,
		reflect.Int32: 32,
		reflect.Int64: 64,
	}
	uintKindToSize = map[reflect.Kind]int{
		reflect.Uint:   0,
		reflect.Uint8:  8,
		reflect.Uint16: 16,
		reflect.Uint32: 32,
		reflect.Uint64: 64,
	}
	floatKindToSize = map[reflect.Kind]int{
		reflect.Float32: 32,
		reflect.Float64: 64,
	}
)

// Decoder is a wrap around csv.Reader
type Decoder struct {
	*csv.Reader
}

// NewDecoder will create a new Decoder to be used
func NewDecoder(r io.Reader) *Decoder {
	dec := &Decoder{csv.NewReader(r)}
	dec.TrimLeadingSpace = true
	return dec
}

// Decode will decode the next line in CSV file to v
func (dec *Decoder) Decode(v interface{}) error {
	var (
		err    error
		record []string
	)
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() || rv.Elem().Kind() != reflect.Struct {
		return errors.New("Decode() expect a pointer to a struct as parameter")
	}

	s := reflect.ValueOf(v).Elem()

	record, err = dec.Read()
	if err != nil {
		return err
	}
	if s.NumField() != len(record) {
		return fmt.Errorf("mismatch length of record: expect %d, get %d", s.NumField(), len(record))
	}

	for i, fValue := range record {
		f := s.Field(i)
		fName := s.Type().Field(i).Name
		if !f.CanSet() {
			return fmt.Errorf("field %q is not settable", fName)
		} else if !f.IsValid() {
			return fmt.Errorf("field %q is not valid", fName)
		}
		// Make sure pointers are properly initialized to nil value
		if f.Kind() == reflect.Ptr && f.IsNil() {
			f.Set(reflect.New(f.Type().Elem()))
		}
		// Only test Unmarshaller interface when it's a pointer and have at
		// least one method.
		if f.Kind() == reflect.Ptr && f.NumMethod() != 0 {
			if x, ok := f.Interface().(Unmarshaller); ok {
				if err = x.UnmarshalCSV(fValue); err != nil {
					return err
				}
				continue
			}
		}
		k := f.Type().Kind()
		if size, ok := intKindToSize[k]; ok {
			var number int64
			if fValue != "" {
				if number, err = strconv.ParseInt(fValue, 10, size); err != nil {
					return fmt.Errorf("failed in parsing %q: %v", fName, err)
				}
			}
			f.SetInt(number)
			continue
		} else if size, ok := uintKindToSize[k]; ok {
			var number uint64
			if fValue != "" {
				if number, err = strconv.ParseUint(fValue, 10, size); err != nil {
					return fmt.Errorf("failed in parsing %q: %v", fName, err)
				}
			}
			f.SetUint(number)
			continue
		} else if size, ok := floatKindToSize[k]; ok {
			var number float64
			if fValue != "" {
				if number, err = strconv.ParseFloat(fValue, size); err != nil {
					return fmt.Errorf("failed in parsing %q: %v", fName, err)
				}
			}
			f.SetFloat(number)
			continue
		} else if k == reflect.String {
			f.SetString(fValue)
			continue
		}
		return fmt.Errorf("don't know how to decode field %q", fName)
	}
	return nil
}
