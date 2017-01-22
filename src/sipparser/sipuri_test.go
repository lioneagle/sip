package sipparser

import (
	"fmt"
	"testing"
)

func BenchmarkSipUri(b *testing.B) {

	for i := 0; i < 10000; i++ {
		fmt.Sprintf("hello")
	}

}

func BenchmarkSipUri2(b *testing.B) {

	for i := 0; i < 1000000; i++ {
		fmt.Sprintf("hello")
	}

}
