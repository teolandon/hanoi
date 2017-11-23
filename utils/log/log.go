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
	f       *os.File
	logger  *log.Logger
	writer  *bufio.Writer
)

func Init() {
	if enabled {
		return
	}

	usr, err := user.Current()
	if err != nil {
		fmt.Println("Logger couldn't get user")
		return
	}

	t := time.Now()

	filename := fmt.Sprintf("hanoi_%04d-%02d-%02d_%02d-%02d-%02d.log",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	err = os.Mkdir(filepath.Join(usr.HomeDir, "logs"), 0777)
	if !os.IsExist(err) {
		return
	}

	f, err = os.Create(filepath.Join(usr.HomeDir, "logs", filename))
	if err != nil {
		fmt.Println("Logger couldn't create file")
		return
	}

	writer = bufio.NewWriter(f)
	logger = log.New(writer, "", log.LstdFlags)
	enabled = true
}

func Log(s ...interface{}) {
	if enabled {
		str := fmt.Sprintln(s)
		logger.Output(2, str)
	}
}

func Close() {
	enabled = false
	writer.Flush()
	f.Close()
}
