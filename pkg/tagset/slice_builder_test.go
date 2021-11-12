package tagset

import (
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSliceBuilder_Fuzz(t *testing.T) {
	f := newNullFactory()
	var lastSb *SliceBuilder
	test := func(seed int64) func(*testing.T) {
		r := rand.New(rand.NewSource(seed))
		return func(t *testing.T) {
			levels := r.Intn(10) + 1
			log.Printf("levels = %d", levels)

			sb := f.NewSliceBuilder(levels, 0)

			// check that we got the same slicebuilder (so it's being reused, behavior
			// we want to validate here)
			if lastSb != nil {
				require.True(t, sb == lastSb, "factory did not reuse SliceBuilder")
			}
			lastSb = sb

			// track the expected values by using an array of builders
			builders := make([]*Builder, levels)
			for l := 0; l < levels; l++ {
				builders[l] = f.NewBuilder(0)
			}

			// add unique tags to the builders
			for tagnum := 0; tagnum < r.Intn(100); tagnum++ {
				level := r.Intn(levels)
				tag := fmt.Sprintf("%d:%d", level, tagnum)
				sb.Add(level, tag)
				builders[level].Add(tag)
			}

			// freeze the builders into tags
			expectedTags := make([]*Tags, levels)
			for l := 0; l < levels; l++ {
				expectedTags[l] = builders[l].Freeze()
				log.Printf("level %d = [%d:%d] = %#v",
					l, sb.offsets[l], sb.offsets[l+1], sb.tags[sb.offsets[l]:sb.offsets[l+1]])
			}

			// verify that each slice is correct (including empty slices)
			for a := 0; a < levels; a++ {
				for b := a; b < levels+1; b++ {
					tags := sb.FreezeSlice(a, b)
					log.Printf("FreezeSlice[%d:%d] = %#v / 0x%016x", a, b, tags.Sorted(), tags.Hash())
					tags.validate(t)

					// union all of the expected tags together
					exp := EmptyTags
					for l := a; l < b; l++ {
						exp = f.Union(exp, expectedTags[l])
					}

					require.Equal(t, exp.Sorted(), tags.Sorted())
				}
			}

			sb.Close()
		}
	}

	// A poor-soul's attempt at fuzzing.  The idea is to catch any edge cases
	// related to reallocation by running a bunch of random, but deterministic
	// (the same on every run), scenarios.
	for i := int64(0); i < 500; i++ {
		t.Run(fmt.Sprintf("seed=%d", i), test(i))
	}
}

func TestSliceBuilder_AddKV(t *testing.T) {
	f := newNullFactory()
	sb := f.NewSliceBuilder(3, 4)

	sb.AddKV(0, "host", "123")
	sb.AddKV(0, "cluster", "k")
	sb.AddKV(1, "container", "abc")
	sb.AddKV(2, "task", "92489")

	t0 := sb.FreezeSlice(0, 1)
	t0.validate(t)
	require.Equal(t, []string{"cluster:k", "host:123"}, t0.Sorted())

	t01 := sb.FreezeSlice(0, 2)
	t01.validate(t)
	require.Equal(t, []string{"cluster:k", "container:abc", "host:123"}, t01.Sorted())

	t012 := sb.FreezeSlice(0, 3)
	t012.validate(t)
	require.Equal(t, []string{"cluster:k", "container:abc", "host:123", "task:92489"}, t012.Sorted())
}
