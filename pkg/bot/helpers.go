package bot

import (
	"fmt"
	"os"
	"strings"
)

func logErr(ff string, args ...interface{}) {
	ff = "[ERROR] " + strings.ReplaceAll(ff, "\n", " ") + "\n"
	fmt.Fprintf(os.Stderr, ff, args...)
}
