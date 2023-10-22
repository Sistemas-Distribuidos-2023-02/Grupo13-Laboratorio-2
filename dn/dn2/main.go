package main

import (
	"context"
	"flag"
    "fmt"
	"log"
    "net"
	"google.golang.org/grpc"
    pb "github.com/VicenteRuizA/proto_lab2/dn_service"
)


var (
	port = flag.Int("port", 50051, "Server port")
)

type saveServer struct {
    pb.UnimplementedSaveServer
}

func (s *saveServer) SaveNaming(ctx context.Context, in *pb.SaveRequest) (*pb.SaveReply, error) {
    log.Printf("Received: \n ID: %v\n Nombre: %v\n Apellido: de %v", in.Id, in.GetName(), in.GetSurname())


	// Use fmt.Sprintf to format the string with variables.
    replyMessage := fmt.Sprintf("Se ha reportado exitosamente ID: %s corresponde a %s %s", in.Id, in.GetName(), in.GetSurname())

	return &pb.SaveReply{Message: replyMessage}, nil
}

func startSaveServer(){
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSaveServer(s, &saveServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	flag.Parse()
    startSaveServer()
}
