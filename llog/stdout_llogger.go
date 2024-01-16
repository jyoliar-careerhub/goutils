package llog

import (
	"encoding/json"
	"fmt"
)

type StdoutLLogger struct{}

func (l *StdoutLLogger) Log(llog *LLog) error {
	bytes, err := json.Marshal(llog)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))
	return nil
}
