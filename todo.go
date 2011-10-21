package main

import (
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

	var err os.Error
	switch flag.Arg(0) {
	case "add":
		t := strings.Join(flag.Args()[1:], " ")
		err = list.AddTask(t)
	case "rm":
		n, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			err = os.NewError("bad argument")
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
