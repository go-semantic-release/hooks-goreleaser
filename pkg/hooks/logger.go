package hooks

import (
	"fmt"
	stdLog "log"
	"os"

	"github.com/apex/log"
)

type LogHandler struct {
	logger *stdLog.Logger
}

func NewLogHandler() *LogHandler {
	return &LogHandler{
		logger: stdLog.New(os.Stderr, "", 0),
	}
}

func (h *LogHandler) HandleLog(e *log.Entry) error {
	str := fmt.Sprintf("(%s) %s", e.Level.String(), e.Message)
	if len(e.Fields) == 0 {
		h.logger.Println(str)
		return nil
	}
	str += " ["
	for k, v := range e.Fields {
		str += fmt.Sprintf("%s=%v ", k, v)
	}
	h.logger.Println(str[0:len(str)-1] + "]")
	return nil
}
