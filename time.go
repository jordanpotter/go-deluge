package deluge

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), `"`)

	if str == "null" {
		t.Time = time.Time{}
		return nil
	}

	unix, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return errors.Wrap(err, "failed to convert to float")
	}

	t.Time = time.Unix(int64(unix), 0)
	return nil
}

func (t *Time) MarshalJSON() ([]byte, error) {
	if t.Time.Unix() == 0 {
		return []byte("null"), nil
	}

	return []byte(string(t.Time.Unix())), nil
}
