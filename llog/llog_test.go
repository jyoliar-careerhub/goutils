package llog

import (
	"testing"
)

func TestLLog(t *testing.T) {
	Default("hello world")
	Log(Msg("hello world"))
	Log(Msg("hello world").Tag("tag1"))
	Log(Msg("hello world").Tags("tag1"))
	Log(Msg("hello world").Tags("tag1", "tag2"))
	Log(Msg("hello world").Tags("tag1", "tag2").Tag("tag3"))
	Log(Msg("Hello world").Data("key1", "value1"))
	Log(Msg("Hello world").Data("key1", "value1").Data("key2", "value2"))
	Log(Msg("Hello world").Datas(map[string]any{"key1": "value1", "key2": "value2", "number": 12}).Data("key3", "value3"))
}
