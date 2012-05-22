package task

import (
	"errors"
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

func (l *List) AddTask(s string, top bool) error {
	var flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	var tasks []string
	if top {
		var err error
		tasks, err = l.Get()
		if err != nil {
			return err
		}
		flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}
	f, err := os.OpenFile(l.filename, flags, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintln(f, s)
	if err != nil {
		return err
	}
	if top {
		for _, t := range tasks {
			_, err = fmt.Fprintln(f, t)
			if err != nil {
				break
			}
		}
	}
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

func (l *List) GetTask(n int) (string, error) {
	tasks, err := l.Get()
	if err != nil {
		return "", err
	}
	if n >= len(tasks) || n < 0 {
		return "", errors.New("index out of range")
	}
	return tasks[n], nil
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
