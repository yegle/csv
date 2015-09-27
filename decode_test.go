package csv_test

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

	"github.com/yegle/csv"
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
		dec := csv.NewDecoder(strings.NewReader(test.input))
		get := reflect.New(reflect.TypeOf(test.expect)).Interface()
		if err := dec.Decode(get); err != nil || !reflect.DeepEqual(reflect.ValueOf(get).Elem().Interface(), test.expect) {
			t.Errorf("Expect %q unmarshal to %#v, get %#v: %v", test.input, test.expect, get, err)
		}
	}
}

func BenchmarkDecodeDefault(b *testing.B) {
	b.ReportAllocs()
	input := `1, 2, 3, 4, 5`
	type T struct {
		F1 int
		F2 int
		F3 int
		F4 int
		F5 int
	}

	t := T{}
	dec := csv.NewDecoder(strings.NewReader(input))
	for i := 0; i < 1000000; i++ {
		dec.Decode(&t)
	}
}

func BenchmarkDecodeUnmarshaller(b *testing.B) {
	b.ReportAllocs()
	input := `"test","test","test","test","test"`
	type T struct {
		F1 *MyString
		F2 *MyString
		F3 *MyString
		F4 *MyString
		F5 *MyString
	}
	t := T{}
	dec := csv.NewDecoder(strings.NewReader(input))
	for i := 0; i < 1000000; i++ {
		dec.Decode(&t)
	}
}

func BenchmarkLendingclubFile(b *testing.B) {
	b.ReportAllocs()
	if TestInput == "" {
		b.Skip("Download loan data from https://www.lendingclub.com/info/download-data.action and specify test file with -input flag")
	}
	var err error
	input, err := os.Open(TestInput)
	if err != nil {
		b.Errorf("failed to open file %q", TestInput)
	}
	defer input.Close()
	stat, err := input.Stat()
	if err != nil {
		b.Errorf("failed to get size of %q", TestInput)
	}
	dec := csv.NewDecoder(input)
	type T struct {
		Fid                      int
		Fmemberid                int
		Floanamnt                int
		Ffundedamnt              int
		Ffundedamntinv           int
		Fterm                    string
		Fintrate                 string
		Finstallment             float32
		Fgrade                   string
		Fsubgrade                string
		Femptitle                string
		Femplength               string
		Fhomeownership           string
		Fannualinc               float32
		Fverificationstatus      string
		Fissued                  *MyDate
		Floanstatus              string
		Fpymntplan               string
		Furl                     string
		Fdesc                    string
		Fpurpose                 string
		Ftitle                   string
		Fzipcode                 string
		Faddrstate               string
		Fdti                     float32
		Fdelinq2yrs              int
		Fearliestcrline          string
		Fficorangelow            int
		Fficorangehigh           int
		Finqlast6mths            int
		Fmthssincelastdelinq     int
		Fmthssincelastrecord     int
		Fopenacc                 int
		Fpubrec                  int
		Frevolbal                int
		Frevolutil               string
		Ftotalacc                int
		Finitialliststatus       string
		Foutprncp                float32
		Foutprncpinv             float32
		Ftotalpymnt              float32
		Ftotalpymntinv           float32
		Ftotalrecprncp           float32
		Ftotalrecint             float32
		Ftotalreclatefee         float32
		Frecoveries              float32
		Fcollectionrecoveryfee   float32
		Flastpymntd              *MyDate
		Flastpymntamnt           float32
		Fnextpymntd              *MyDate
		Flastcreditpulld         *MyDate
		Flastficorangehigh       int
		Flastficorangelow        int
		Fcollections12mthsexmed  int
		Fmthssincelastmajorderog int
		Fpolicycode              string
	}

	t := T{}
	for {
		err = dec.Decode(&t)
		if err == io.EOF {
			break
		} else if err != nil {
			b.Log(err)
		}
	}
	b.SetBytes(stat.Size())
}
