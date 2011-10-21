package task

import (
	"bufio"
	"fmt"
	"os"
)

type List struct {
	filename string
}

func NewList(filename string) *List {
	return &List{filename}
}

func (l *List) AddTask(s string) os.Error {
	const flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	f, err := os.OpenFile(l.filename, flags, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintln(f, s)
	return err
}

func (l *List) RemoveTask(n int) os.Error {
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

func (l *List) Get() ([]string, os.Error) {
	f, err := os.Open(l.filename)
	if err != nil {
		if err, ok := err.(*os.PathError); ok && err.Error == os.ENOENT {
			// A non-existent file means no tasks.
			return nil, nil
		}
		return nil, err
	}
	var tasks []string
	br := bufio.NewReader(f)
	for {
		// TODO(adg): handle long lines if prefix bool is set
		t, _, err := br.ReadLine()
		if err == os.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, string(t))
	}
	return tasks, nil
}
