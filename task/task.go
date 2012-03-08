package task

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type List struct {
	filename string
}

func NewList(filename string) *List {
	return &List{filename}
}

func (l *List) AddTask(s string) error {
	const flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	f, err := os.OpenFile(l.filename, flags, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintln(f, s)
	return err
}

func (l *List) RemoveTask(n int) error {
	tasks, err := l.Get()
	if err != nil {
		return err
	}
	f, err := os.Create(l.filename)
	if err != nil {
		return err
	}
	defer f.Close()
	for i, t := range tasks {
		if i == n {
			continue
		}
		_, err = fmt.Fprintln(f, t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *List) Get() ([]string, error) {
	f, err := os.Open(l.filename)
	if err != nil {
		if os.IsNotExist(err) {
			// A non-existent file means no tasks.
			return nil, nil
		}
		return nil, err
	}
	var tasks []string
	br := bufio.NewReader(f)
	for {
		t, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, string(t))
	}
	return tasks, nil
}
