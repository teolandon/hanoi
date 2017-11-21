package log

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

var (
	enabled bool = false
	logger  *log.Logger
)

func init() {
	usr, err := user.Current()
	if err != nil {
		return
	}

	t := time.Now()

	filename := fmt.Sprintf("%4d-%2d-%2d_%2d-%2d-%2d.log", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	f, err := os.Create(filepath.Join(usr.HomeDir, "logs", filename))
	if err != nil {
		return
	}

	defer f.Close()

	writer := bufio.NewWriter(f)
	logger = log.New(writer, "", log.LstdFlags|log.Lshortfile)
	enabled = true
}

func Log(s string) {
	logger.Println(s)
}
