package fault

import (
	"fmt"
	"runtime"
)

type Fault struct {
	Message  string `json:"message"`
	funcName string `json:"-"`
	fileName string `json:"-"`
}

func (f *Fault) Error() string {
	return f.Message
}

func (f *Fault) String() string {
	return fmt.Sprintf("[%s:%s] %s", f.fileName, f.funcName, f.Message)
}

func New(msg string) *Fault {
	pc, filename, line, _ := runtime.Caller(1)

	return &Fault{
		Message:  msg,
		funcName: runtime.FuncForPC(pc).Name(),
		fileName: fmt.Sprintf("%s:%d", filename, line),
	}
}

func (f *Fault) Newf(format string, args ...any) *Fault {
	pc, filename, line, _ := runtime.Caller(1)

	return &Fault{
		Message:  fmt.Sprintf(format, args...),
		funcName: runtime.FuncForPC(pc).Name(),
		fileName: fmt.Sprintf("%s:%d", filename, line),
	}
}

func (f *Fault) MarshalJSON() ([]byte, error) {
	return []byte(f.String()), nil
}

func (f *Fault) UnmarshalJSON(data []byte) error {
	f.Message = string(data)
	return nil
}
