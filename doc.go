/*
Package csv is a wrap around csv package in standard library. It provides a way
to unmarshal lines in CSV file to a struct.

Decoder embed a csv.Reader so all csv.Reader's property can be used on Decoder.

Only Decoder is available right now. Encoder has lower priority so maybe when I
have time...

Example code:

	CSVText := `1, 2, 3, "test", My string
	4, 5, 6, "another_test", my_other_string`

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
*/
package csv
