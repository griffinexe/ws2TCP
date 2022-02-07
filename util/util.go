package util

import "io"

func CopyWorker(src io.Reader, dst io.Writer, doneCh chan bool) {
	io.Copy(dst, src)
	doneCh <- true
}
