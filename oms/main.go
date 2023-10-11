package main

import (
	"context"
	"flag"
    "fmt"
	"log"
    "net"
	"google.golang.org/grpc"
    pb "github.com/VicenteRuizA/proto_lab2/in/continente_post_oms"
)


var (
	port = flag.Int("port", 50051, "Server port")
)

type server struct {
    pb.UnimplementedReportServer
}

func (s *server) IdentifyCondition(ctx context.Context, in *pb.SeverityRequest) (*pb.SeverityReply, error) {
    log.Printf("Received: \n Nombre: %v\n Apellido: %v\nCondicion de %v", in.GetName(), in.GetSurname(), in.GetCondition())

	// Use fmt.Sprintf to format the string with variables.
	replyMessage := fmt.Sprintf("Se ha reportado exitosamente que %s %s esta %s", in.GetName(), in.GetSurname(), in.GetCondition())

	return &pb.SeverityReply{Message: replyMessage}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterReportServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
