package llog

import (
	"fmt"
	"strings"
	"time"
)

type StdoutLLogger struct {
}

func (l *StdoutLLogger) Log(llog *LLog) error {
	msg := msgString(llog.msg)
	tags := tagsString(llog.tags)
	datas := datasString(llog.datas)
	createdAt := createdAtString(llog.createdAt)

	fmt.Printf("%s%s%s%s\n", createdAt, msg, tags, datas)

	return nil
}
func msgString(msg string) string {
	if msg == "" {
		return ""
	}

	return fmt.Sprintf(" %s", msg)
}

func tagsString(tags []string) string {
	if len(tags) == 0 {
		return ""
	}

	return fmt.Sprintf(" [%s]", strings.Join(tags, ","))
}

func datasString(datas map[string]string) string {
	if len(datas) == 0 {
		return ""
	}

	var datasString []string
	for key, value := range datas {
		datasString = append(datasString, fmt.Sprintf("%s=%s", key, value))
	}

	return fmt.Sprintf(" {%s}", strings.Join(datasString, ","))
}

func createdAtString(createdAt time.Time) string {
	return createdAt.Format("2006-01-02 15:04:05")
}
