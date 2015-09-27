package csv

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

var TestInput string

func init() {
	flag.StringVar(&TestInput, "input", "", "Benchmark test input file")
	flag.Parse()
}

type MyString string

func (s *MyString) UnmarshalCSV(data string) error {
	*s = MyString(fmt.Sprintf("%10s", data))
	return nil
}

func NewMyString(s string) *MyString {
	return (*MyString)(&s)
}

type MySimpleString string

func (s *MySimpleString) UnmarshalCSV(data string) error {
	*s = MySimpleString(data)
	return nil
}

type MyDate struct {
	time.Time
}

func (d *MyDate) UnmarshalCSV(data string) error {
	if data == "" {
		return nil
	}
	parsed, err := time.Parse("Jan-2006", data)
	if err != nil {
		return err
	}
	*d = MyDate{parsed}
	return nil
}

func TestDecode(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	tests := []struct {
		input  string
		expect interface{}
	}{
		{
			input: `1, 2,"test", 5`,
			expect: struct {
				F1 int
				F2 uint
				F3 string
				F4 float64
			}{1, 2, "test", 5.0},
		},
		{
			input: `1, 2, 3, 4, 5`,
			expect: struct {
				F1 int
				F2 int
				F3 int
				F4 int
				F5 int
			}{1, 2, 3, 4, 5},
		},
		{
			input: `1, 2, "string", 5`,
			expect: struct {
				F1 int
				F2 int
				F3 *MyString
				F4 int
			}{
				1, 2, NewMyString("    string"), 5,
			},
		},
	}
	for _, test := range tests {
		dec := NewDecoder(strings.NewReader(test.input))
		get := reflect.New(reflect.TypeOf(test.expect)).Interface()
		if err := dec.Decode(get); err != nil || !reflect.DeepEqual(reflect.ValueOf(get).Elem().Interface(), test.expect) {
			t.Errorf("Expect %q unmarshal to %#v, get %#v: %v", test.input, test.expect, get, err)
		}
	}
}

func BenchmarkDecodeDefault(b *testing.B) {
	input := `1, 2, 3, 4, 5`
	type T struct {
		F1 int
		F2 int
		F3 int
		F4 int
		F5 int
	}

	t := T{}
	dec := NewDecoder(strings.NewReader(input))
	for i := 0; i < 1000000; i++ {
		dec.Decode(&t)
	}
}

func BenchmarkDecodeUnmarshaller(b *testing.B) {
	input := `"test","test","test","test","test"`
	type T struct {
		F1 MySimpleString
		F2 MySimpleString
		F3 MySimpleString
		F4 MySimpleString
		F5 MySimpleString
	}
	t := T{}
	dec := NewDecoder(strings.NewReader(input))
	for i := 0; i < 1000000; i++ {
		dec.Decode(&t)
	}
}

func BenchmarkLendingclubFile(b *testing.B) {
	if TestInput == "" {
		return
	}
	input, err := os.Open(TestInput)
	if err != nil {
		b.Errorf("failed to open file %q", TestInput)
	}
	defer input.Close()
	stat, err := input.Stat()
	if err != nil {
		b.Errorf("failed to get size of %q", TestInput)
	}
	dec := NewDecoder(input)
	type T struct {
		Fid                          int
		Fmember_id                   int
		Floan_amnt                   int
		Ffunded_amnt                 int
		Ffunded_amnt_inv             int
		Fterm                        string
		Fint_rate                    string
		Finstallment                 float32
		Fgrade                       string
		Fsub_grade                   string
		Femp_title                   string
		Femp_length                  string
		Fhome_ownership              string
		Fannual_inc                  float32
		Fverification_status         string
		Fissue_d                     *MyDate
		Floan_status                 string
		Fpymnt_plan                  string
		Furl                         string
		Fdesc                        string
		Fpurpose                     string
		Ftitle                       string
		Fzip_code                    string
		Faddr_state                  string
		Fdti                         float32
		Fdelinq_2yrs                 int
		Fearliest_cr_line            string
		Ffico_range_low              int
		Ffico_range_high             int
		Finq_last_6mths              int
		Fmths_since_last_delinq      int
		Fmths_since_last_record      int
		Fopen_acc                    int
		Fpub_rec                     int
		Frevol_bal                   int
		Frevol_util                  string
		Ftotal_acc                   int
		Finitial_list_status         string
		Fout_prncp                   float32
		Fout_prncp_inv               float32
		Ftotal_pymnt                 float32
		Ftotal_pymnt_inv             float32
		Ftotal_rec_prncp             float32
		Ftotal_rec_int               float32
		Ftotal_rec_late_fee          float32
		Frecoveries                  float32
		Fcollection_recovery_fee     float32
		Flast_pymnt_d                *MyDate
		Flast_pymnt_amnt             float32
		Fnext_pymnt_d                *MyDate
		Flast_credit_pull_d          *MyDate
		Flast_fico_range_high        int
		Flast_fico_range_low         int
		Fcollections_12_mths_ex_med  int
		Fmths_since_last_major_derog int
		Fpolicy_code                 string
	}

	t := T{}
	for {
		err := dec.Decode(&t)
		if err == io.EOF {
			break
		} else if err != nil {
			b.Log(err)
		}
	}
	b.SetBytes(stat.Size())
}
