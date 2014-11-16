// +build test_zproto_c
package example

// #cgo CFLAGS: -I../../../../include -I/usr/local/include
// #cgo LDFLAGS: ../../../.libs/libzproto.a
// #cgo pkg-config: libczmq
//
// #include <czmq.h>
// #include <zproto_example.h>
//
// static zsock_t* zproto_zsock_new(int type) {
//     return zsock_new(type);
// }
//
// static int zproto_zsock_bind(zsock_t *sock, const char *addr) {
//     return zsock_bind(sock, addr);
// }
//
// static void zproto_zsock_destroy(zsock_t **sock) {
//     zsock_destroy(sock);
// }
import "C"

import (
	"testing"
	"unsafe"

	zmq "github.com/pebbe/zmq4"
)

// Tests interoperability of Go implementation with C implementation
func testLogInputFromC(t *testing.T) {

	// Output socket (C implementaion)
	output := C.zproto_zsock_new(C.ZMQ_DEALER)
	C.zproto_zsock_bind(output, C.CString("tcp://127.0.0.1:5557"))

	// Input socket
	input, err := zmq.NewSocket(zmq.ROUTER)
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()

	err = input.Connect("tcp://127.0.0.1:5557")
	if err != nil {
		t.Fatal(err)
	}
	defer input.Disconnect("tcp://127.0.0.1:5557")

	// Create a Log message and send it through the wire
	log := C.zproto_example_new(C.ZPROTO_EXAMPLE_LOG)
	C.zproto_example_set_sequence(log, 120)
	C.zproto_example_set_level(log, 121)
	C.zproto_example_set_event(log, 122)
	C.zproto_example_set_node(log, 123)
	C.zproto_example_set_peer(log, 124)
	C.zproto_example_set_time(log, 125)
	//log.Host = "Life is short but Now lasts for ever"
	//log.Data = "Life is short but Now lasts for ever"

	// Send twice from same object
	C.zproto_example_send_again(log, unsafe.Pointer(output))
	C.zproto_example_send(&log, unsafe.Pointer(output))

	for i := 0; i < 2; i++ {
		transit, err := Recv(input)
		if err != nil {
			t.Fatal(err)
		}

		tr := transit.(*Log)

		if tr.sequence != 120 {
			t.Fatalf("expected %d, got %d", 123, tr.sequence)
		}

		if tr.Level != 121 {
			t.Fatalf("expected %d, got %d", 123, tr.Level)
		}

		if tr.Event != 122 {
			t.Fatalf("expected %d, got %d", 123, tr.Event)
		}

		if tr.Node != 123 {
			t.Fatalf("expected %d, got %d", 123, tr.Node)
		}

		if tr.Peer != 124 {
			t.Fatalf("expected %d, got %d", 123, tr.Peer)
		}

		if tr.Time != 125 {
			t.Fatalf("expected %d, got %d", 123, tr.Time)
		}

		if tr.Host != "Life is short but Now lasts for ever" {
			//t.Fatalf("expected %s, got %s", "Life is short but Now lasts for ever", tr.Host)
		}

		if tr.Data != "Life is short but Now lasts for ever" {
			//t.Fatalf("expected %s, got %s", "Life is short but Now lasts for ever", tr.Data)
		}
	}

	/*
		err = tr.Send(input)
		if err != nil {
			t.Fatal(err)
		}

		transit, err = Recv(output)
		if err != nil {
			t.Fatal(err)
		}

		if routingId != string(tr.RoutingId()) {
			t.Fatalf("expected %s, got %s", routingId, string(tr.RoutingId()))
		}
	*/
	C.zproto_zsock_destroy(&output)
}
