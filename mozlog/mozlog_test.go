package mozlog

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMozLogger(t *testing.T) {
	in := new(bytes.Buffer)
	DefaultLogger.Output = in
	UseMozLogger("testlogger")
	log.Println("test")
	assert.Contains(t, in.String(), "mozlog_test")
}
