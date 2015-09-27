/*
Package csv is a wrap around csv package in standard library. It provides a way
to unmarshal lines in CSV file to a struct.

Decoder embed a csv.Reader so all csv.Reader's property can be used on Decoder.

Example code:

	type MyString string
	func (s *MyString) UnmarshalCSV(data string) error {
		*s = data
		return nil
	}
	type T struct {
		F1 int
		F2 int
		F3 float
		F4 string
		F5 *MyString
	}
	input, _ := os.Open("input.csv")
	dec := csv.NewDecoder(input)
	dec.TrimLeadingSpace = true

	t := T{}
	dec.Decode(&t)
	fmt.Println(t)
*/
package csv
