package util

import "io"

func copyWorker(src io.Reader, dst io.Writer, doneCh chan bool) {
	io.Copy(dst, src)
	doneCh <- true
}

func IOCopy(io1, io2 io.ReadWriteCloser) {
	ch := make(chan bool)
	go copyWorker(io1, io2, ch)
	go copyWorker(io2, io1, ch)
	<-ch
	io1.Close()
	io2.Close()
	<-ch
}
