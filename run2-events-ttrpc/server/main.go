package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/containerd/containerd/api/services/ttrpc/events/v1"
	ttrpc "github.com/containerd/ttrpc"

	apievents "github.com/containerd/containerd/api/events"

	"github.com/containerd/typeurl"
	"github.com/gogo/protobuf/types"
)

const socket = "run2-events.ttrpc"

type server struct{}

func (s *server) Forward(ctx context.Context, req *pb.ForwardRequest) (*types.Empty, error) {
	fmt.Println("-------------- Event ----------------")
	fmt.Printf("Topic: %s\n", req.Envelope.Topic)
	fmt.Printf("Namespace: %s\n", req.Envelope.Namespace)
	fmt.Printf("Time: %s\n", req.Envelope.Timestamp)

	e, err := typeurl.UnmarshalAny(req.Envelope.Event)
	if err != nil {
		log.Fatalf("Can not unmarshal event: %v", err)
	}

	switch e.(type) {
	case *apievents.TaskCreate:
		te, _ := e.(*apievents.TaskCreate)
		fmt.Printf("Event: %v", te)

	case *apievents.TaskStart:
		te, _ := e.(*apievents.TaskStart)
		fmt.Printf("Event: %v", te)

	case *apievents.TaskOOM:
		te, _ := e.(*apievents.TaskOOM)
		fmt.Printf("Event: %v", te)

	case *apievents.TaskExit:
		te, _ := e.(*apievents.TaskExit)
		fmt.Printf("ContainerID: %s\n", te.ContainerID)
		fmt.Printf("ID: %s\n", te.ID)
		fmt.Printf("Pid: %d\n", te.Pid)
		fmt.Printf("ExitStatus: %d\n", te.ExitStatus)
		fmt.Printf("ExitedAt: %s\n", te.ExitedAt)

	case *apievents.TaskDelete:
		te, _ := e.(*apievents.TaskDelete)
		fmt.Printf("Event: %v", te)

	case *apievents.TaskExecAdded:
		te, _ := e.(*apievents.TaskExecAdded)
		fmt.Printf("Event: %v", te)

	case *apievents.TaskExecStarted:
		te, _ := e.(*apievents.TaskExecStarted)
		fmt.Printf("Event: %v", te)

	case *apievents.TaskPaused:
		te, _ := e.(*apievents.TaskPaused)
		fmt.Printf("Event: %v", te)

	case *apievents.TaskResumed:
		te, _ := e.(*apievents.TaskResumed)
		fmt.Printf("Event: %v", te)

	case *apievents.TaskCheckpointed:
		te, _ := e.(*apievents.TaskCheckpointed)
		fmt.Printf("Event: %v", te)
	default:
		log.Fatal("Event type not found")
	}

	return &types.Empty{}, nil
}

func main() {
	s, err := ttrpc.NewServer()
	if err != nil {
		log.Fatalf("failed creating ttrpc server: %v", err)
	}
	defer s.Close()

	pb.RegisterEventsService(s, &server{})

	if err := os.Remove(socket); err != nil {
		log.Printf("Failed removing socket: %v", err)
	}

	l, err := net.Listen("unix", socket)
	if err != nil {
		log.Fatalf("failed creating listener: %v", err)
	}
	defer l.Close()

	fmt.Println("Running...")

	if err := s.Serve(context.Background(), l); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
