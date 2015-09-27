/*
Package csv is a wrap around csv package in standard library. It provides a way
to unmarshal lines in CSV file to a struct.

Decoder embed a csv.Reader so all csv.Reader's property can be used on Decoder.

Only Decoder is available right now. Encoder has lower priority so maybe when I
have time...
*/
package csv
