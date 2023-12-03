package sample

import "fmt"

func SampleError(s string) error {
	return fmt.Errorf("sample error %s", s)
}
