package changroup_test

import (
	"testing"

	"github.com/frantjc/forge/changroup"
)

func TestAllSettled(t *testing.T) {
	var (
		aC = make(chan any, 1)
		bC = make(chan any, 1)
		i  = 0
	)

	close(aC)

	t.Run("", func(t *testing.T) {
		if i != 0 {
			t.Error("was", i, "but expected", 0)
			t.FailNow()
		}

		t.Parallel()

		<-changroup.AllSettled(aC, bC)

		if i != 1 {
			t.Error("was", i, "but expected", 1)
			t.FailNow()
		}
	})

	i++

	close(bC)
}
