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

type saveServer struct {
	pb.UnimplementedSaveServer
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

func LeerPrimeraLinea() {
	file, err := os.Open("DATA.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "-")
		if len(fields) == 3 {
			fmt.Printf("ID: %s, Nombre: %s, Apellido: %s\n", fields[0], fields[1], fields[2])
		} else {
			fmt.Println("El formato de la línea no es válido")
		}
	}
}

func (s *saveServer) SaveNaming(ctx context.Context, in *pb.SaveRequest) (*pb.SaveReply, error) {
	log.Printf("Received: \n ID: %v\n Nombre: %v\n Apellido: %v", in.Id, in.GetName(), in.GetSurname())
	writeToDataFile(in.Id, in.GetName(), in.GetSurname())
	LeerPrimeraLinea()
	// Use fmt.Sprintf to format the string with variables.
	replyMessage := fmt.Sprintf("Se ha reportado exitosamente ID: %s corresponde a %s %s", in.Id, in.GetName(), in.GetSurname())
	return &pb.SaveReply{Message: replyMessage}, nil
}

func startSaveServer() {
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
