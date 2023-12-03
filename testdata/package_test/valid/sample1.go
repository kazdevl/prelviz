package valid

import (
	"fmt"

	"github.com/kazdevl/prelviz/testdata/package_test/valid/nest/sample"
)

func Sample1() {
	fmt.Printf("error %v", sample.SampleError("sample1"))
}
