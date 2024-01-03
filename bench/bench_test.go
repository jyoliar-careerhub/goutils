package bench

import (
	"crypto/rand"
	"testing"
)

type FieldPtr struct {
	A *string
	B *string
	C *string
	D *string
}

type FieldValue struct {
	A string
	B string
	C string
	D string
}

func NewPtr4FieldPtr(a, b, c, d *string) *FieldPtr {

	return &FieldPtr{
		A: a,
		B: b,
		C: b,
		D: b,
	}
}

func NewPtr4FieldValue(a, b, c, d string) *FieldValue {
	return &FieldValue{
		A: a,
		B: b,
		C: c,
		D: d,
	}
}

func NewValue4FieldPtr(a, b, c, d *string) FieldPtr {

	return FieldPtr{
		A: a,
		B: b,
		C: c,
		D: d,
	}
}

func NewValue4FieldValue(a, b, c, d string) FieldValue {

	return FieldValue{
		A: a,
		B: b,
		C: c,
		D: d,
	}
}

func BenchmarkTestPointerNValue(b *testing.B) {

	b.Run("NewPtr4FieldPtr", func(b *testing.B) {
		randStrs := generateRandom4String(5000)
		for i := 0; i < b.N; i++ {
			sample := NewPtr4FieldPtr(&randStrs[0], &randStrs[1], &randStrs[2], &randStrs[3])

			if *(*sample).A == "" || *(*sample).B == "" || *(*sample).C == "" || *(*sample).D == "" {
				panic("empty string")
			}
		}
	})

	b.Run("NewPtr4FieldValue", func(b *testing.B) {
		randStrs := generateRandom4String(5000)

		for i := 0; i < b.N; i++ {
			sample := NewPtr4FieldValue(randStrs[0], randStrs[1], randStrs[2], randStrs[3])
			if sample.A == "" || sample.B == "" || sample.C == "" || sample.D == "" {
				panic("empty string")
			}
		}
	})

	b.Run("NewValue4FieldPtr", func(b *testing.B) {
		randStrs := generateRandom4String(5000)

		for i := 0; i < b.N; i++ {
			sample := NewValue4FieldPtr(&randStrs[0], &randStrs[1], &randStrs[2], &randStrs[3])

			if *sample.A == "" || *sample.B == "" || *sample.C == "" || *sample.D == "" {
				panic("empty string")
			}
		}
	})

	b.Run("NewValue4FieldValue", func(b *testing.B) {
		randStrs := generateRandom4String(5000)

		for i := 0; i < b.N; i++ {
			sample := NewValue4FieldValue(randStrs[0], randStrs[1], randStrs[2], randStrs[3])

			if sample.A == "" || sample.B == "" || sample.C == "" || sample.D == "" {
				panic("empty string")
			}
		}
	})

}

func CheckEmptyValue(fv FieldValue) bool {
	return fv.A == "" || fv.B == "" || fv.C == "" || fv.D == ""
}

func CheckEmptyValueViaPtr(fv *FieldValue) bool {
	return fv.A == "" || fv.B == "" || fv.C == "" || fv.D == ""
}

func BenchmarkFieldValue(b *testing.B) {
	randStrs := generateRandom4String(30000)
	sample := NewPtr4FieldValue(randStrs[0], randStrs[1], randStrs[2], randStrs[3])

	b.Run("CheckEmptyValue", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			CheckEmptyValue(*sample)
		}
	})

	b.Run("CheckEmptyValueViaPtr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			CheckEmptyValueViaPtr(sample)
		}
	})
}

func BenchmarkChannel(b *testing.B) {
	randStrs := generateRandom4String(3000)
	sample := NewPtr4FieldValue(randStrs[0], randStrs[1], randStrs[2], randStrs[3])

	b.Run("PtrChannel", func(b *testing.B) {
		channel := make(chan *FieldValue)

		for i := 0; i < b.N; i++ {
			go func() {
				channel <- sample
			}()
			recvSample := <-channel

			CheckEmptyValueViaPtr(recvSample)
		}

		close(channel)
	})

	b.Run("ValueChannel", func(b *testing.B) {
		channel := make(chan FieldValue)

		for i := 0; i < b.N; i++ {
			go func() {
				channel <- *sample
			}()
			recvSample := <-channel

			CheckEmptyValue(recvSample)
		}

		close(channel)
	})
}

func generateRandom4String(length int) []string {
	return []string{
		generateRandomString(length),
		generateRandomString(length),
		generateRandomString(length),
		generateRandomString(length),
	}
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	for i := 0; i < length; i++ {
		randomBytes[i] = charset[randomBytes[i]%byte(len(charset))]
	}

	return string(randomBytes)
}
