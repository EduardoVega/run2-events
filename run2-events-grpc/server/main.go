package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	cdevents "github.com/containerd/containerd/api/events"
	pb "github.com/containerd/containerd/api/services/events/v1"
	"github.com/containerd/typeurl"
	"github.com/gogo/protobuf/types"
	//epb "github.com/EduardoVega/run2-events-grpc/events"
)

type server struct {
	pb.UnimplementedEventsServer
}

func (*server) Forward(ctx context.Context, req *pb.ForwardRequest) (*types.Empty, error) {
	fmt.Println("#################################")
	fmt.Printf("Topic: %s\n", req.Envelope.Topic)
	fmt.Printf("Namespace: %s\n", req.Envelope.Namespace)
	fmt.Printf("Time: %s\n", req.Envelope.Timestamp)

	e, err := typeurl.UnmarshalAny(req.Envelope.Event)
	if err != nil {
		log.Fatalf("Can not unmarshal event: %v", err)
	}

	switch e.(type) {
	case *cdevents.TaskCreate:
		te, _ := e.(*cdevents.TaskCreate)
		fmt.Printf("Event: %v", te)

	case *cdevents.TaskStart:
		te, _ := e.(*cdevents.TaskStart)
		fmt.Printf("Event: %v", te)

	case *cdevents.TaskOOM:
		te, _ := e.(*cdevents.TaskOOM)
		fmt.Printf("Event: %v", te)

	case *cdevents.TaskExit:
		te, _ := e.(*cdevents.TaskExit)
		fmt.Printf("ContainerID: %s\n", te.ContainerID)
		fmt.Printf("ID: %s\n", te.ID)
		fmt.Printf("Pid: %d\n", te.Pid)
		fmt.Printf("ExitStatus: %d\n", te.ExitStatus)
		fmt.Printf("ExitedAt: %s\n", te.ExitedAt)

	case *cdevents.TaskDelete:
		te, _ := e.(*cdevents.TaskDelete)
		fmt.Printf("Event: %v", te)

	case *cdevents.TaskExecAdded:
		te, _ := e.(*cdevents.TaskExecAdded)
		fmt.Printf("Event: %v", te)

	case *cdevents.TaskExecStarted:
		te, _ := e.(*cdevents.TaskExecStarted)
		fmt.Printf("Event: %v", te)

	case *cdevents.TaskPaused:
		te, _ := e.(*cdevents.TaskPaused)
		fmt.Printf("Event: %v", te)

	case *cdevents.TaskResumed:
		te, _ := e.(*cdevents.TaskResumed)
		fmt.Printf("Event: %v", te)

	case *cdevents.TaskCheckpointed:
		te, _ := e.(*cdevents.TaskCheckpointed)
		fmt.Printf("Event: %v", te)
	default:
		log.Fatal("Event type not found")
	}

	fmt.Println("#################################")

	return &types.Empty{}, nil
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterEventsServer(s, &server{})

	if err := s.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
