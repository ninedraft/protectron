package registry

import (
	"time"

	"github.com/asdine/storm"
	bolt "github.com/coreos/bbolt"
)

func defaultStormOptions() func(*storm.Options) error {
	return storm.BoltOptions(0600, // file mode
		&bolt.Options{
			Timeout: 1 * time.Second,
		})
}
