package movingavg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type mockItem struct {
	time.Time
}

func (m mockItem) GetTime() time.Time {
	return m.Time
}

func newMockItem(t time.Time) Item {
	return mockItem{t}
}

type mockNower struct {
	now time.Time
}

func (nower mockNower) Now() time.Time {
	return nower.now

}

func TestCountFrom(t *testing.T) {
	type testcase struct {
		list     []Item
		expected int
	}

	now := time.Now()

	testcases := []testcase{
		{
			[]Item{newMockItem(now)},
			1,
		},
		{
			[]Item{newMockItem(now), newMockItem(now)},
			2,
		},
		{
			[]Item{
				newMockItem(now.Add(time.Duration(-3) * time.Second)),
				newMockItem(now),
			},
			1,
		}, {
			[]Item{
				newMockItem(now.Add(time.Duration(-3) * time.Second)),
				newMockItem(now.Add(time.Duration(-2) * time.Second)),
				newMockItem(now.Add(time.Duration(-1) * time.Second)),
				newMockItem(now),
			},
			3,
		},
	}

	for _, testcase := range testcases {
		cal := NewCalculator(mockNower{now: now})
		for _, item := range testcase.list {
			cal.AddItem(item)
		}

		got := cal.CountFrom(time.Duration(2) * time.Second)

		wanted := testcase.expected

		assert.Equal(t, wanted, got, fmt.Sprintf("l: %#v w: %#v g: %#v", testcase.list, wanted, got))
	}
}
