package main

import (
	"context"
	"fmt"
	"log"
	"time"

	cdevents "github.com/containerd/containerd/api/events"
	pb "github.com/containerd/containerd/api/services/events/v1"

	// epb "github.com/EduardoVega/run2-events-grpc/events"
	"github.com/containerd/typeurl"
	"google.golang.org/grpc"
)

func main() {
	con, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed connecting: %v", err)
	}
	defer con.Close()

	client := pb.NewEventsClient(con)

	event := &cdevents.TaskExit{
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

	res, err := client.Forward(context.Background(), req)
	if err != nil {
		log.Fatalf("Response error: %v", err)
	}

	fmt.Printf("Response: %v", res.String())
}
