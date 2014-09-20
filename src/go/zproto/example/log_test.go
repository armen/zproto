package example

import (
	"testing"

	zmq "github.com/pebbe/zmq4"
)

// Yay! Test function.
func TestLog(t *testing.T) {

	// Create pair of sockets we can send through

	// Output socket
	output, err := zmq.NewSocket(zmq.DEALER)
	if err != nil {
		t.Fatal(err)
	}
	defer output.Close()

	routingId := "Shout"
	output.SetIdentity(routingId)
	err = output.Bind("inproc://selftest-log")
	if err != nil {
		t.Fatal(err)
	}
	defer output.Unbind("inproc://selftest-log")

	// Input socket
	input, err := zmq.NewSocket(zmq.ROUTER)
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()

	err = input.Connect("inproc://selftest-log")
	if err != nil {
		t.Fatal(err)
	}
	defer input.Disconnect("inproc://selftest-log")

	// Create a Log message and send it through the wire
	log := NewLog()
	log.sequence = 120
	log.Level = 121
	log.Event = 122
	log.Node = 123
	log.Peer = 124
	log.Time = 125
	log.Host = "Life is short but Now lasts for ever"
	log.Data = "Life is short but Now lasts for ever"

	err = log.Send(output)
	if err != nil {
		t.Fatal(err)
	}
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
		t.Fatalf("expected %s, got %s", "Life is short but Now lasts for ever", tr.Host)
	}

	if tr.Data != "Life is short but Now lasts for ever" {
		t.Fatalf("expected %s, got %s", "Life is short but Now lasts for ever", tr.Data)
	}

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
}
