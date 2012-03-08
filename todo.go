package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/nf/todo/task"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	defaultFile := filepath.Join(os.Getenv("HOME"), ".todo")
	file := flag.String("file", defaultFile, "file to store tasks")
	flag.Parse()

	list := task.NewList(*file)

	var err error
	switch flag.Arg(0) {
	case "add":
		t := strings.Join(flag.Args()[1:], " ")
		err = list.AddTask(t)
	case "rm":
		var n int
		n, err = strconv.Atoi(flag.Arg(1))
		if err != nil {
			err = errors.New("bad task number")
			break
		}
		err = list.RemoveTask(n)
	default:
		var tasks []string
		tasks, err = list.Get()
		for i, t := range tasks {
			fmt.Printf("%2d: %s\n", i, t)
		}
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
