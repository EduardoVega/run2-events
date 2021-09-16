package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/containerd/containerd/api/services/ttrpc/events/v1"
	"github.com/containerd/ttrpc"

	apievents "github.com/containerd/containerd/api/events"
	"github.com/containerd/typeurl"
)

const socket = "run2-events.ttrpc"

func main() {
	conn, err := net.Dial("unix", socket)
	if err != nil {
		log.Fatalf("failed connecting to server: %v", err)
	}
	defer conn.Close()

	tc := ttrpc.NewClient(conn)

	client := pb.NewEventsClient(tc)

	event := &apievents.TaskExit{
		ContainerID: "1",
		ID:          "1",
		Pid:         1,
		ExitStatus:  1,
		ExitedAt:    time.Now(),
	}

	any, err := typeurl.MarshalAny(event)
	if err != nil {
		log.Fatalf("Error marshaling event: %v", err)
	}

	req := &pb.ForwardRequest{
		Envelope: &pb.Envelope{
			Timestamp: time.Now(),
			Namespace: "default",
			Topic:     "/tasks/exit",
			Event:     any,
		},
	}

	_, err = client.Forward(context.Background(), req)
	if err != nil {
		log.Fatalf("Error forwarding request: %v", err)
	}
}
