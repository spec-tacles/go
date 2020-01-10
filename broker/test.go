package broker

type testReadWriter struct {
	C chan []byte
}

func (r *testReadWriter) Read(d []byte) (int, error) {
	b := <-r.C
	return copy(d, b), nil
}

func (r *testReadWriter) Write(d []byte) (int, error) {
	r.C <- d
	return len(d), nil
}

var rcv = make(chan *IOPacket)
var cb = func(event string, data []byte) {
	rcv <- &IOPacket{event, data}
}
