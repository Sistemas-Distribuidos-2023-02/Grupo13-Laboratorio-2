package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	pb "github.com/VicenteRuizA/proto_lab2/dn_service"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "Server port")
)

type server struct {
	pb.UnimplementedSaveServer
	pb.UnimplementedLoadServer
}

func writeToDataFile(id string, nombre string, apellido string) error {
	file, err := os.OpenFile("DATA.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	line := fmt.Sprintf("%s-%s-%s\n", id, nombre, apellido)

	_, err = file.WriteString(line)
	if err != nil {
		return err
	}

	return nil
}

func LeerNombreApellido(id_in string) (string, string) {

	file, _ := os.Open("DATA.txt")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "-")
		if len(fields) == 3 {
			if fields[0] == id_in {
				name := fields[1]
				surname := fields[2]

				return name, surname
			}
		}
	}
	return "error", "error"
}

func (s *server) SaveNaming(ctx context.Context, in *pb.SaveRequest) (*pb.SaveReply, error) {
	log.Printf("Received: \n ID: %v\n Nombre: %v\n Apellido: %v", in.Id, in.GetName(), in.GetSurname())
	writeToDataFile(in.Id, in.GetName(), in.GetSurname())
	replyMessage := fmt.Sprintf("Se ha reportado exitosamente ID: %s corresponde a %s %s", in.Id, in.GetName(), in.GetSurname())
	return &pb.SaveReply{Message: replyMessage}, nil
}

func (s *server) RequestData(ctx context.Context, in *pb.DataRequest) (*pb.DataReply, error) {
	log.Printf("Received: \n ID: %v", in.Id)

	nombre, apellido := LeerNombreApellido(in.Id)

	return &pb.DataReply{Nombre: nombre, Apellido: apellido}, nil
}

func startServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSaveServer(s, &server{}) // register Save service
	pb.RegisterLoadServer(s, &server{}) // register Load service
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	flag.Parse()
	startServer()
}
