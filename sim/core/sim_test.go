package core

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

var lines = make([]string, 0, 10e6)

var numLines = 0

func BenchmarkLog(b *testing.B) {
	message := "%s took %d damage from %s"
	vals := []interface{}{"something", 433, "sinister strike"}

	dur := 3456 * time.Millisecond

	b.Run("fmt + fmt", func(b *testing.B) {
		b.ReportAllocs()
		lines = lines[:0]
		numLines = 0
		for i := 0; i < b.N; i++ {
			lines = append(lines, fmt.Sprintf("[%0.2f] "+fmt.Sprintf(message, vals...)+"\n", dur.Seconds()))
		}
		numLines += len(lines)
	})
	b.Logf("sample - %s\n", lines[0])
	b.Logf("fmt + fmt lines = %d\n", numLines)

	b.Run("fmt append", func(b *testing.B) {
		b.ReportAllocs()
		lines = lines[:0]
		numLines = 0
		for i := 0; i < b.N; i++ {
			lines = append(lines, fmt.Sprintf("[%0.2f] "+message+"\n", append([]interface{}{dur.Seconds()}, vals...)...))
		}
		numLines += len(lines)
	})
	b.Logf("sample - %s\n", lines[0])
	b.Logf("fmt append lines = %d\n", numLines)

	b.Run("fmt strconv", func(b *testing.B) {
		b.ReportAllocs()
		lines = lines[:0]
		numLines = 0
		for i := 0; i < b.N; i++ {
			lines = append(lines, "["+strconv.FormatFloat(dur.Seconds(), 'f', 2, 64)+"] "+fmt.Sprintf(message, vals...)+"\n")
		}
		numLines += len(lines)
	})
	b.Logf("sample - %s\n", lines[0])
	b.Logf("fmt strconv lines = %d\n", numLines)

	b.Run("fmt fmt", func(b *testing.B) {
		b.ReportAllocs()
		lines = lines[:0]
		numLines = 0
		for i := 0; i < b.N; i++ {
			lines = append(lines, fmt.Sprintf("[%0.2f] %s\n", dur.Seconds(), fmt.Sprintf(message, vals...)))
		}
		numLines += len(lines)
	})
	b.Logf("sample - %s\n", lines[0])
	b.Logf("fmt fmt lines = %d\n", numLines)
}
