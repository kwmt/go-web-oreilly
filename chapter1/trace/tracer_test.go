package trace
import (
	"testing"
	"bytes"
)

func TestNew(t *testing.T) {


	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Newからの戻り値がnilです")
	} else {
		tracer.Trace("こんにちは、trace パッケージ")
		if buf.String() != "こんにちは、trace パッケージ\n" {
			t.Errorf("'%s'という誤った文字列が出力されました。", buf.String())
		}
	}

}