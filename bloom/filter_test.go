package bloom_test

import (
	"crypto/sha256"
	"testing"

	"github.com/AugustineAurelius/yuki/bloom"
	"github.com/stretchr/testify/require"
)

func Test_Filter(t *testing.T) {
	filter := bloom.NewFilter(bloom.NewBitSet(), 1_000_000, 64)
	data := []byte("test data")
	filter.Add(data)

	require.True(t, filter.Test(data))
	require.False(t, filter.Test([]byte("test dato")))
}

func Benchmark_Filter(b *testing.B) {
	filter := bloom.NewFilter(bloom.NewBitSet(), 1_000_000, 64)
	hash := sha256.New()
	hash.Write([]byte("test data"))
	data := hash.Sum(nil)

	for i := 0; i < b.N; i++ {
		filter.Add(data)
	}

}
