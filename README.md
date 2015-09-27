# csv

[![GoDoc](https://godoc.org/github.com/yegle/csv?status.svg)](https://godoc.org/github.com/yegle/csv)

`csv` allows you unmarshal and marshal (TBD) between csv (Comma
Separated Values) text and Golang struct.

# Quick Example

```go
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
```
