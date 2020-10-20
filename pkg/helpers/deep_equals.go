package helpers

import (
	"fmt"
	"time"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/check.v1"
)

var timeComparer = cmp.Comparer(func(x, y time.Time) bool {
	return x.Unix() == y.Unix()
})

type deepEqualsChecker struct {
	*check.CheckerInfo
}

// DeepEquals checker verifies that the obtained value is deep-equal to
// the expected value. but compared to the "check" version this will also
// compare time.Time type.
var DeepEquals check.Checker = &deepEqualsChecker{
	&check.CheckerInfo{Name: "DeepEquals", Params: []string{"obtained", "expected"}},
}

func (checker *deepEqualsChecker) Check(params []interface{}, names []string) (result bool, error string) {
	result = cmp.Equal(params[1], params[0], timeComparer)
	if !result {
		error = fmt.Sprintf("mismatch (-want +got):\n%s", cmp.Diff(params[1], params[0], timeComparer))
	}
	return
}
