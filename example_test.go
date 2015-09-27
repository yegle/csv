package csv_test

import (
	"fmt"
	"io"
	"strings"

	"github.com/yegle/csv"
)

func ExampleDecoder_Decode() {
	CSVText := `1, 2, 3, "test"
	4, 5, 6, "another_test"`

	type T struct {
		F1 int
		F2 int
		F3 float32
		F4 string
	}
	dec := csv.NewDecoder(strings.NewReader(CSVText))
	dec.TrimLeadingSpace = true

	t := T{}
	for {
		err := dec.Decode(&t)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error: %v\n", err)
		}
		fmt.Println(t)
	}
	//Output:
	//{1 2 3 test}
	//{4 5 6 another_test}
}
