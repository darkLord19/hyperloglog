package hyperloglog

import (
	"errors"
	"hash"
	"hash/fnv"
	"math"
	"math/bits"
	"math/rand"
	"strconv"
)

type HyperLogLog struct {
	HashFunction hash.Hash64
	Accuracy     float64
	store        []int
	size         int
	b            int
}

const (
	DEFAULT_ACCURACY = 70
)

func getIndexingBitsSize(accuracy float64) uint {
	return uint(
		math.Ceil(
			math.Log2(
				math.Pow((1.04 / ((100 - accuracy) / 100)), 2),
			),
		),
	)
}

func New(options ...func(*HyperLogLog)) (*HyperLogLog, error) {
	hll := &HyperLogLog{}

	for _, option := range options {
		option(hll)
	}

	if hll.Accuracy == 0 {
		hll.Accuracy = DEFAULT_ACCURACY
	}
	if hll.HashFunction == nil {
		hll.HashFunction = fnv.New64()
	}

	if hll.Accuracy >= 100 || hll.Accuracy <= 0 {
		return nil, errors.New("accuracy must be between 0 and 100")
	}

	hll.b = getIndexingBitsSize(hll.Accuracy)
	hll.size = uint(math.Exp2(float64(hll.b)))
	hll.store = make([]uint8, hll.size)
	return hll, nil
}

// Option methods
func WithAccuracy(accuracy float64) func(*HyperLogLog) {
	return func(hll *HyperLogLog) {
		hll.Accuracy = accuracy
	}
}

func (b *HyperLogLog) getHash(seed int, key []byte) (uint64, error) {
	b.HashFunction.Reset()
	t := []byte(strconv.Itoa(seed))
	var err error
	_, err = b.HashFunction.Write(t)
	_, err = b.HashFunction.Write(key)
	return b.HashFunction.Sum64(), err
}

func WithHash(hash hash.Hash64) func(*HyperLogLog) {
	return func(hll *HyperLogLog) {
		hll.HashFunction = hash
	}
}

func (hll *HyperLogLog) Add(element []byte) error {
	hash, err := hll.getHash(rand.Int(), element)
	if err != nil {
		return err
	}
	l := uint64((1<<hll.b)-1) & hash
	r := hash >> (64 - hll.b)
	leadingZeroes := bits.LeadingZeros64(r)
	hll.store[l] = int(math.Max(float64(leadingZeroes), float64(hll.store[l])))
	return nil
}
