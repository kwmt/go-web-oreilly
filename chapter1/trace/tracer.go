package trace
import (
	"io"
	"fmt"
)

// コード内での出来事を記録できるオブジェクトを表すインターフェース
type Tracer interface {
	Trace(...interface{})
}


func New(w io.Writer) Tracer{
	// 公開されないtracerを返す意味は、
	// ユーザーは単に「Tracerインターフェースに合致したオブジェクト」を受け取るだけで、 privateなtracer型には関知しない。
	// ユーザーはインターフェースに基いて操作を行う。
	// そのため、tracerが他のメソッドやフィール土を公開していてもユーザーからは見えず、パッケージのAPIをクリーンデシンプルな状態に保てる。
	
	return &tracer{out: w}
}


type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}){
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}
