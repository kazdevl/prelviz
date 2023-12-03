package valid

import (
	"fmt"
	"time"

	"github.com/kazdevl/prelviz/testdata/package_test/valid/nest/sample"
)

func Sample() {
	v := sample.NewSample(time.Now())
	fmt.Println(v.Name)
	fmt.Println(sample.SampleString())
}
