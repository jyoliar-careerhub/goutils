package llog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLLog(t *testing.T) {

	checkLog(t, func() error { Fatal("hello world"); return nil }, Msg("hello world").Level(FATAL))
	checkLog(t, func() error { Error("hello world"); return nil }, Msg("hello world").Level(ERROR))
	checkLog(t, func() error { Warn("hello world"); return nil }, Msg("hello world").Level(WARN))
	checkLog(t, func() error { Info("hello world"); return nil }, Msg("hello world").Level(INFO))
	checkLog(t, func() error { Debug("hello world"); return nil }, Msg("hello world").Level(DEBUG))

	checkLogFormat(t, Msg("hello world").Tag("tag1"))
	checkLogFormat(t, Msg("hello world").Tags("tag1"))
	checkLogFormat(t, Msg("hello world").Tags("tag1", "tag2"))
	checkLogFormat(t, Msg("hello world").Tags("tag1", "tag2").Tag("tag3"))
	checkLogFormat(t, Msg("Hello world").Data("key1", "value1"))
	checkLogFormat(t, Msg("Hello world").Data("key1", "value1").Data("key2", "value2"))
	checkLogFormat(t, Msg("Hello world").Datas(map[string]any{"key1": "value1", "key2": "value2", "isTrue": true, "isFalse": false}).Data("key3", "value3"))
	checkLogFormat(t, Msg("Hello world").Datas(map[string]any{"key1": "value1", "key2": "value2", "number": 12}).Data("key3", "value3"))
}

func checkLogFormat(t *testing.T, expected *LLogBuilder) {

	logBuf, err := interceptStdout(func() error { return expected.Log() })
	require.NoError(t, err)
	now := time.Now()

	var actualLog LLog
	err = json.Unmarshal(logBuf.Bytes(), &actualLog)
	require.NoError(t, err)

	require.Equal(t, time.Time(actualLog.CreatedAt).Unix(), now.Unix())

	assertEqualLogs(t, expected.Build(), &actualLog)
}

func checkLog(t *testing.T, action func() error, expected *LLogBuilder) {
	logBuf, err := interceptStdout(action)
	require.NoError(t, err)

	var actualLog LLog
	err = json.Unmarshal(logBuf.Bytes(), &actualLog)
	require.NoError(t, err)

	assertEqualLogs(t, expected.Build(), &actualLog)
}

func assertEqualLogs(t *testing.T, expected *LLog, actual *LLog) {
	expectedCreated := expected.CreatedAt
	actualCreated := actual.CreatedAt
	defer func() {
		expected.CreatedAt = expectedCreated
		actual.CreatedAt = actualCreated
	}()

	// ignore time
	expected.CreatedAt = LogTime{}
	actual.CreatedAt = LogTime{}

	for actualK, actualV := range actual.Datas {
		// require.Equal(t, expected.Datas[actualK], actualV)
		require.Equal(t, fmt.Sprintf("%v", expected.Datas[actualK]), fmt.Sprintf("%v", actualV))
	}
	require.Len(t, actual.Datas, len(expected.Datas))

	for _, acactualTag := range actual.Tags {
		require.Contains(t, expected.Tags, acactualTag)
	}
	require.Len(t, actual.Tags, len(expected.Tags))
}

func interceptStdout(action func() error) (bytes.Buffer, error) {
	originStdout := os.Stdout
	defer func() { os.Stdout = originStdout }()

	r, w, _ := os.Pipe()
	os.Stdout = w

	err := action()
	if err != nil {
		return bytes.Buffer{}, err
	}

	w.Close()

	outC := make(chan bytes.Buffer)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf
	}()

	return <-outC, nil
}
