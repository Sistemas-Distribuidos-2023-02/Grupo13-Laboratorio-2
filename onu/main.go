package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "github.com/VicenteRuizA/proto_lab2/oms_service"
)

const (
	defaultcondition = "INFECTADO"
)

var (

    // onu envia request a oms 
    // oms se encuentra en vm049
    addr = flag.String("addr", "10.6.46.59:50051", "ip address to connect to")

    //probar local
    //addr = flag.String("addr", "localhost:50051", "ip address to connect to")
	condition = flag.String("condition", defaultcondition, "Condition to request")
)


func main() {
	flag.Parse()
	
    conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		log.Fatalf("fallo la conexion: %v", err)
	}
	defer conn.Close()

    c := pb.NewRequestClient(conn)


	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
    r, err := c.RequestCondition(ctx, &pb.ConditionRequest{Condition : *condition})
    if err != nil{
        log.Fatalf("fallo en request: %v", err)
    }	
    log.Printf("respuesta: %s %s", r.GetNombre(), r.GetApellido())
	
}
