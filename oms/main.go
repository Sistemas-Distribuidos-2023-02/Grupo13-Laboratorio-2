package main

import (
	"context"
	"flag"
    "fmt"
	"log"
    "net"
    "time"
    //"sync"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
    pbs "github.com/VicenteRuizA/proto_lab2/in/continente_post_oms"
    pbc "github.com/VicenteRuizA/proto_lab2/in/oms_save_dn"
)
// notar que pbs es protocol buffer server, dado que oms es servidor en dicha communicacion


// server -------------------------------------------------------------------------------------------

var (
	port = flag.Int("port", 50051, "Server port")
)


type server struct {
    pbs.UnimplementedReportServer
}

func (s *server) IdentifyCondition(ctx context.Context, in *pbs.SeverityRequest) (*pbs.SeverityReply, error) {
    log.Printf("Received: \n Nombre: %v\n Apellido: %v\n Condicion de %v", in.GetName(), in.GetSurname(), in.GetCondition())



    // Generar id aquí?
    // modularizar lo de abajo en una función?


    name_id := "1"
    
    //probar en en vms donde oms, servidor al cual pide requiest, esta en 049
    addr :=  "10.6.46.59:50051"

    //probar en en vms donde oms, servidor al cual pide requiest, esta en 049
    //addr :=  "localhost:50052"
  
    conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		log.Fatalf("fallo la conexion: %v", err)
	}
	defer conn.Close()

	c := pbc.NewSaveClient(conn)


	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	

    r, err := c.SaveNaming(ctx, &pbc.SaveRequest{Id : name_id, Name : in.GetName(), Surname : in.GetSurname()})
    if err != nil{
        log.Fatalf("could not greet: %v", err)
    }	
    log.Printf("Reply from server: %s", r.GetMessage())

	replyMessage := fmt.Sprintf("Se ha reportado exitosamente que %s %s esta %s", in.GetName(), in.GetSurname(), in.GetCondition())

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
