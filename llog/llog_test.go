package llog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLLog(t *testing.T) {

	t.Run("Default Log functions", func(t *testing.T) {
		ctx := context.TODO()
		t.Run("Fatal", func(t *testing.T) {
			checkStdoutLog(t, func() error { return Fatal(ctx, "hello world") }, Msg("hello world").Level(FATAL))
		})

		t.Run("Error", func(t *testing.T) {
			checkStdoutLog(t, func() error { return Error(ctx, "hello world") }, Msg("hello world").Level(ERROR))
		})

		t.Run("Warn", func(t *testing.T) {
			checkStdoutLog(t, func() error { return Warn(ctx, "hello world") }, Msg("hello world").Level(WARN))
		})

		t.Run("Info", func(t *testing.T) {
			checkStdoutLog(t, func() error { return Info(ctx, "hello world") }, Msg("hello world").Level(INFO))
		})

		t.Run("Debug", func(t *testing.T) {
			checkStdoutLog(t, func() error { return Debug(ctx, "hello world") }, Msg("hello world").Level(DEBUG))
		})
	})

	t.Run("Logging with builder", func(t *testing.T) {

		t.Run("Msg_Tag", func(t *testing.T) {
			checkStdoutLogFormat(t, Msg("hello world").Tag("tag1"))
		})

		t.Run("Msg_one Tags", func(t *testing.T) {
			checkStdoutLogFormat(t, Msg("hello world").Tags("tag1"))
		})

		t.Run("Msg_two Tags", func(t *testing.T) {
			checkStdoutLogFormat(t, Msg("hello world").Tags("tag1", "tag2"))
		})

		t.Run("Msg_Tags_tag", func(t *testing.T) {
			checkStdoutLogFormat(t, Msg("hello world").Tags("tag1", "tag2").Tag("tag3"))
		})

		t.Run("Msg_Data", func(t *testing.T) {
			checkStdoutLogFormat(t, Msg("Hello world").Data("key1", "value1"))
		})

		t.Run("Msg_Data_Data", func(t *testing.T) {
			checkStdoutLogFormat(t, Msg("Hello world").Data("key1", "value1").Data("key2", "value2"))
		})

		t.Run("Msg_DatasWithBool", func(t *testing.T) {
			checkStdoutLogFormat(t, Msg("Hello world").Datas(map[string]any{"key1": "value1", "key2": "value2", "isTrue": true, "isFalse": false}).Data("key3", "value3"))
		})

		t.Run("Msg_DatasWithNumber_Data", func(t *testing.T) {
			checkStdoutLogFormat(t, Msg("Hello world").Datas(map[string]any{"key1": "value1", "key2": "value2", "number": 12}).Data("key3", "value3"))
		})
	})

	t.Run("Logging with error", func(t *testing.T) {
		t.Run("not LLogError", func(t *testing.T) {
			checkStdoutLogErr(t, fmt.Errorf("hello world"), Msg("hello world").Level(ERROR))
		})

		t.Run("LLogError_NoTags_NoDatas", func(t *testing.T) {
			checkStdoutLogErr(t, &LLogError{Err: fmt.Errorf("hello world")}, Msg("hello world").Level(ERROR))
		})

		t.Run("LLogError_NoTags_Data", func(t *testing.T) {
			checkStdoutLogErr(t, &LLogError{Err: fmt.Errorf("hello world"), datas: map[string]any{"key1": "value1"}}, Msg("hello world").Level(ERROR).Data("key1", "value1"))
			checkStdoutLogErr(t, NewLLogError(fmt.Errorf("hello world")).Data("key1", "value1"), Msg("hello world").Level(ERROR).Data("key1", "value1"))
			checkStdoutLogErr(t, NewLLogError(fmt.Errorf("hello world")).Data("bool", true).Data("string", "value2").Data("number", 12), Msg("hello world").Level(ERROR).Datas(map[string]any{"bool": true, "string": "value2", "number": 12}))

		})

		t.Run("LLogError_NoTags_Datas", func(t *testing.T) {
			checkStdoutLogErr(t, &LLogError{Err: fmt.Errorf("hello world"), datas: map[string]any{"bool": true, "string": "value2", "number": 12}}, Msg("hello world").Level(ERROR).Datas(map[string]any{"bool": true, "string": "value2", "number": 12}))
			checkStdoutLogErr(t, NewLLogError(fmt.Errorf("hello world")).Datas(map[string]any{"bool": true, "string": "value2", "number": 12}), Msg("hello world").Level(ERROR).Datas(map[string]any{"bool": true, "string": "value2", "number": 12}))
			checkStdoutLogErr(t, NewLLogError(fmt.Errorf("hello world")).Datas(map[string]any{"bool": true, "string": "value2"}).Data("number", 12), Msg("hello world").Level(ERROR).Datas(map[string]any{"bool": true, "string": "value2", "number": 12}))
		})

		t.Run("LLogError_Tag_NoDatas", func(t *testing.T) {
			checkStdoutLogErr(t, &LLogError{Err: fmt.Errorf("hello world"), tags: []string{"tag1"}}, Msg("hello world").Level(ERROR).Tag("tag1"))
			checkStdoutLogErr(t, NewLLogError(fmt.Errorf("hello world")).Tag("tag1"), Msg("hello world").Level(ERROR).Tag("tag1"))
			checkStdoutLogErr(t, NewLLogError(fmt.Errorf("hello world")).Tag("tag1").Tag("tag2"), Msg("hello world").Level(ERROR).Tags("tag1", "tag2"))
			checkStdoutLogErr(t, NewLLogError(fmt.Errorf("hello world")).Tag("tag1").Tag("tag2").Tag("tag3"), Msg("hello world").Level(ERROR).Tag("tag1").Tag("tag2").Tag("tag3"))
		})

		t.Run("LLogError_Tags_NoDatas", func(t *testing.T) {
			checkStdoutLogErr(t, &LLogError{Err: fmt.Errorf("hello world"), tags: []string{"tag1", "tag2"}}, Msg("hello world").Level(ERROR).Tags("tag1", "tag2"))
			checkStdoutLogErr(t, NewLLogError(fmt.Errorf("hello world")).Tags("tag1", "tag2"), Msg("hello world").Level(ERROR).Tags("tag1", "tag2"))
			checkStdoutLogErr(t, NewLLogError(fmt.Errorf("hello world")).Tags("tag1", "tag2").Tag("tag3"), Msg("hello world").Level(ERROR).Tags("tag1", "tag2", "tag3"))
			checkStdoutLogErr(t, NewLLogError(fmt.Errorf("hello world")).Tags("tag1", "tag2").Tag("tag3").Tag("tag4"), Msg("hello world").Level(ERROR).Tags("tag1", "tag2", "tag3", "tag4"))
		})
	})
	resetDefault := func() {
		SetDefaultDatas(make(map[string]any))
		SetDefaultTags([]string{})
		SetDefaultContextDatas([]string{})
	}
	t.Run("Logging with default datas", func(t *testing.T) {
		ctx := context.TODO()

		t.Run("DefaultDatas", func(t *testing.T) {
			resetDefault()
			SetDefaultDatas(map[string]any{"key1": "value1", "key2": "value2"})
			// checkStdoutLog(t, func() error { return Msg("hello world").Level(INFO).Data("key3", "value3").Log() }, Msg("hello world").Level(INFO).Datas(map[string]any{"key1": "value1", "key2": "value2", "key3": "value3"}))
			// checkStdoutLog(t, func() error { return Msg("hello world").Level(INFO).Data("key2", "value2_modified").Log() }, Msg("hello world").Level(INFO).Datas(map[string]any{"key1": "value1", "key2": "value2_modified"}))
			checkStdoutLog(t, func() error { return Debug(ctx, "hello world") }, Msg("hello world").Level(DEBUG).Datas(map[string]any{"key1": "value1", "key2": "value2"}))
		})

		t.Run("DefaultData", func(t *testing.T) {
			resetDefault()
			SetDefaultData("key1", "value1")
			checkStdoutLog(t, func() error { return Msg("hello world").Level(INFO).Data("key2", "value2").Log(ctx) }, Msg("hello world").Level(INFO).Datas(map[string]any{"key1": "value1", "key2": "value2"}))
			checkStdoutLog(t, func() error { return Msg("hello world").Level(INFO).Data("key1", "value1_modified").Log(ctx) }, Msg("hello world").Level(INFO).Datas(map[string]any{"key1": "value1_modified"}))
			checkStdoutLog(t, func() error { return Debug(ctx, "hello world") }, Msg("hello world").Level(DEBUG).Datas(map[string]any{"key1": "value1"}))
		})

		t.Run("DefaultTags", func(t *testing.T) {
			resetDefault()
			SetDefaultTags([]string{"tag1", "tag2"})
			checkStdoutLog(t, func() error { return Msg("hello world").Level(INFO).Tag("tag3").Log(ctx) }, Msg("hello world").Level(INFO).Tags("tag1", "tag2", "tag3"))
			checkStdoutLog(t, func() error { return Msg("hello world").Level(INFO).Tag("tag2").Log(ctx) }, Msg("hello world").Level(INFO).Tags("tag1", "tag2"))
			checkStdoutLog(t, func() error { return Debug(ctx, "hello world") }, Msg("hello world").Level(DEBUG).Tags("tag1", "tag2"))
		})

		t.Run("DefaultTag", func(t *testing.T) {
			resetDefault()
			SetDefaultTag("tag1")
			checkStdoutLog(t, func() error { return Msg("hello world").Level(INFO).Tag("tag2").Log(ctx) }, Msg("hello world").Level(INFO).Tags("tag1", "tag2"))
			checkStdoutLog(t, func() error { return Msg("hello world").Level(INFO).Tag("tag1").Log(ctx) }, Msg("hello world").Level(INFO).Tags("tag1"))
			checkStdoutLog(t, func() error { return Debug(ctx, "hello world") }, Msg("hello world").Level(DEBUG).Tags("tag1"))
		})
	})

	t.Run("Logging with context datas", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "string", "value1")
		ctx = context.WithValue(ctx, "bool", true)
		ctx = context.WithValue(ctx, "number", 12)

		t.Run("data is not existed", func(t *testing.T) {
			resetDefault()
			SetDefaultContextData("NOT_EXISTED")
			checkStdoutLog(t, func() error { return Msg("hello world").Level(INFO).Log(ctx) }, Msg("hello world").Level(INFO))
		})

		t.Run("data is existed", func(t *testing.T) {
			t.Run("string", func(t *testing.T) {
				resetDefault()
				SetDefaultContextData("string")
				checkStdoutLog(t, func() error { return Msg("hello world").Level(INFO).Log(ctx) }, Msg("hello world").Level(INFO).Data("string", "value1"))
			})

			t.Run("bool", func(t *testing.T) {
				resetDefault()
				SetDefaultContextData("bool")
				checkStdoutLog(t, func() error { return Msg("hello world").Level(INFO).Log(ctx) }, Msg("hello world").Level(INFO).Data("bool", true))
			})

			t.Run("number", func(t *testing.T) {
				resetDefault()
				SetDefaultContextData("number")
				checkStdoutLog(t, func() error { return Msg("hello world").Level(INFO).Log(ctx) }, Msg("hello world").Level(INFO).Data("number", 12))
			})

			t.Run("string_bool_number", func(t *testing.T) {
				resetDefault()
				SetDefaultContextDatas([]string{"string", "bool", "number"})
				checkStdoutLog(t, func() error { return Msg("hello world").Level(INFO).Log(ctx) }, Msg("hello world").Level(INFO).Datas(map[string]any{"string": "value1", "bool": true, "number": 12}))
			})

			t.Run("override context data", func(t *testing.T) {
				resetDefault()
				SetDefaultContextData("string")
				checkStdoutLog(t, func() error { return Msg("hello world").Level(INFO).Data("string", "value2_modified").Log(ctx) }, Msg("hello world").Level(INFO).Data("string", "value2_modified"))
			})
		})
	})
}

func checkStdoutLogErr(t *testing.T, err error, expected *LLogBuilder) {
	ctx := context.TODO()
	logBuf, err := interceptStdout(func() error { return LogErr(ctx, err) })
	require.NoError(t, err)

	var actualLog LLog
	err = json.Unmarshal(logBuf.Bytes(), &actualLog)
	require.NoError(t, err)

	assertEqualLogs(t, expected.Build(), &actualLog)
}

func checkStdoutLogFormat(t *testing.T, expected *LLogBuilder) {
	ctx := context.TODO()
	logBuf, err := interceptStdout(func() error { return expected.Log(ctx) })
	require.NoError(t, err)
	now := time.Now()

	var actualLog LLog
	err = json.Unmarshal(logBuf.Bytes(), &actualLog)
	require.NoError(t, err)

	require.Equal(t, time.Time(actualLog.CreatedAt).Unix(), now.Unix())

	assertEqualLogs(t, expected.Build(), &actualLog)
}

func checkStdoutLog(t *testing.T, action func() error, expected *LLogBuilder) {
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
