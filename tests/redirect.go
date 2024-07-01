package tests

import (
	"bytes"
	"io"
	"os"
)

type RedirectOutput struct {
	oldStdout *os.File
	outputCh  chan string
}

func RedirectStdoutToChannel() (*os.File, *RedirectOutput) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outputCh := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outputCh <- buf.String()
	}()

	return w, &RedirectOutput{
		oldStdout: oldStdout,
		outputCh:  outputCh,
	}
}

func (ro *RedirectOutput) RedirectChannelToStdout(w *os.File) string {
	w.Close()
	os.Stdout = ro.oldStdout
	output := <-ro.outputCh
	return output
}
