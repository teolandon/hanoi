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
		fmt.Println("Couldn't get user")
		return
	}

	t := time.Now()

	filename := fmt.Sprintf("%4d-%2d-%2d_%2d-%2d-%2d.log", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	err = os.Mkdir(filepath.Join(usr.HomeDir, "logs"), 0777)
	if !os.IsExist(err) {
		return
	}

	f, err := os.Create(filepath.Join(usr.HomeDir, "logs", filename))
	if err != nil {
		fmt.Println("Couldn't create file")
		return
	}

	writer := bufio.NewWriter(f)
	logger = log.New(writer, "", log.LstdFlags|log.Lshortfile)
	enabled = true
}

func Log(s ...interface{}) {
	if enabled {
		logger.Output(2, fmt.Sprintln(s))
	}
}
