package jjson

import (
	"encoding/json"
	"testing"
)

type Sample struct {
	Name string
	Age  int
}

func BenchmarkJjson(b *testing.B) {
	b.Run("jjson.Unmarshal", func(b *testing.B) {
		data := []byte(`{"Name":"John","Age":30}`)
		for i := 0; i < b.N; i++ {
			sample, err := Unmarshal[Sample](data)
			if err != nil {
				b.Fatal(err)
			}
			_ = sample
		}
	})

	b.Run("jjson.UnmarshalStack", func(b *testing.B) {
		data := []byte(`{"Name":"John","Age":30}`)
		for i := 0; i < b.N; i++ {
			sample, err := Unmarshal2(data, &Sample{})
			if err != nil {
				b.Fatal(err)
			}
			_ = sample
		}
	})

	b.Run("json.Unmarshal", func(b *testing.B) {
		data := []byte(`{"Name":"John","Age":30}`)
		for i := 0; i < b.N; i++ {
			var sample Sample
			err := json.Unmarshal(data, &sample)
			if err != nil {
				b.Fatal(err)
			}
			_ = sample
		}
	})
}
