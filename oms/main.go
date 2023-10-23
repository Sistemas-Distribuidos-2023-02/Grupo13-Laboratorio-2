package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	//"sync"
	pbc "github.com/VicenteRuizA/proto_lab2/dn_service"
	pbs "github.com/VicenteRuizA/proto_lab2/oms_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// notar que pbs es protocol buffer server, dado que oms es servidor en dicha communicacion

var (
	port    = flag.Int("port", 50051, "Server port")
	name_id = 1
)

// escribe la oms en DATA.txt (modificable para los datanodes)
func writeToDataFile(id int, datanote string, estado string) error {
	file, err := os.OpenFile("DATA.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	line := fmt.Sprintf("%d-datanode%s-%s\n", id, datanote, estado)

	_, err = file.WriteString(line)
	if err != nil {
		return err
	}

	return nil
}

func conexionADatanode(name string, condition string, name_id int) string {
	if name[0] >= 'A' && name[0] <= 'M' {
		fmt.Println("El primer caracter está en el rango A-M")
		err := writeToDataFile(name_id, "1", condition)
		if err != nil {
			log.Fatal(err)
		}
		return "10.6.46.61:50051"
		// se envia el mensaje (id, nombre y apellido)

	} else if name[0] >= 'N' && name[0] <= 'Z' {
		fmt.Println("El primer caracter está en el rango N-Z")
		err := writeToDataFile(name_id, "2", condition)
		if err != nil {
			log.Fatal(err)
		}
		return "10.6.46.62:50051"
	}
	return "error"
}

type server struct {
	pbs.UnimplementedReportServer
}

func (s *server) IdentifyCondition(ctx context.Context, in *pbs.SeverityRequest) (*pbs.SeverityReply, error) {
	log.Printf("Received: \n Nombre: %v\n Apellido: %v\n Condicion de %v", in.GetName(), in.GetSurname(), in.GetCondition())
	// Generar id aquí?

	// modularizar lo de abajo en una función?

	// oms es cliente de dn
	// notar dn1 en vm051
	// dn2 vm052
	//addr_dn1 := "10.6.46.61:50051"
	//addr_dn2 :=  "10.6.46.62:50051"

	//probar en en vms donde oms, servidor al cual pide requiest, esta en 049
	//addr :=  "localhost:50052"

	addr := conexionADatanode(in.GetName(), in.GetCondition(), name_id)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fallo la conexion: %v", err)
	}
	defer conn.Close()

	c := pbc.NewSaveClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, err := c.SaveNaming(ctx, &pbc.SaveRequest{Id: strconv.Itoa(name_id), Name: in.GetName(), Surname: in.GetSurname()})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Reply from server: %s", r.GetMessage())

	replyMessage := fmt.Sprintf("Se ha reportado exitosamente que %s %s esta %s", in.GetName(), in.GetSurname(), in.GetCondition())
	name_id += 1

	return &pbs.SeverityReply{Message: replyMessage}, nil
}

func startServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pbs.RegisterReportServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// main -------------------------------------------------------------------------------------------
func main() {
	flag.Parse()
	startServer()
}
