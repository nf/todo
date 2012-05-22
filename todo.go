package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nf/todo/task"
)

var noAct = errors.New("didn't act")

const usage = `Usage:
	todo
		Show top task
	todo ls
		Show all tasks
	todo N
		Promote task N to top
	todo rm
		Remove top task
	todo rm N
		Remove task N
	todo description...
		Add task with given description
Flags:
`

func main() {
	defaultFile := filepath.Join(os.Getenv("HOME"), ".todo")
	if f := os.Getenv("TODO"); f != "" {
		defaultFile = f
	}
	file := flag.String("file", defaultFile, "file in which to store tasks")
	now := flag.Bool("now", false, "when adding, insert at head")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
		flag.PrintDefaults()
	}
	flag.Parse()

	list := task.NewList(*file)
	a, n := flag.Arg(0), len(flag.Args())

	err := noAct
	switch {
	case a == "ls" && n == 1:
		// list tasks
		var tasks []string
		tasks, err = list.Get()
		for i := len(tasks) - 1; i >= 0; i-- {
			fmt.Printf("%2d: %s\n", i, tasks[i])
		}
	case a == "rm" && n <= 2:
		// remove [latest] task
		i, err2 := strconv.Atoi(flag.Arg(1))
		if err2 != nil && n == 2 {
			break
		}
		if n == 1 {
			i = 0
		}
		err = list.RemoveTask(i)
		if err != nil || n == 2 {
			break
		}
		fallthrough
	case n == 0:
		// no arguments; show top task
		var tasks []string
		tasks, err = list.Get()
		if len(tasks) > 0 {
			fmt.Println(tasks[0])
		}
	case n == 1:
		// single argument; might be shift-to-top
		i, err2 := strconv.Atoi(flag.Arg(0))
		if err2 != nil {
			// not a number, probably just a single word todo
			break
		}
		var tasks []string
		tasks, err = list.Get()
		if err != nil {
			break
		}
		if i >= len(tasks) || i < 0 {
			err = errors.New("index out of range")
		}
		t := tasks[i]
		err = list.RemoveTask(i)
		if err != nil {
			break
		}
		err = list.AddTask(t, true)
		if err == nil {
			fmt.Println(t)
		}
	}
	if err == noAct {
		// no action taken, assume add
		t := strings.Join(flag.Args(), " ")
		err = list.AddTask(t, *now)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
