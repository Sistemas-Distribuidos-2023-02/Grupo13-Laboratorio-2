package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/VicenteRuizA/proto_lab2/oms_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// const (
// 	defaultcondition = "INFECTADO"
// )

var (

	// onu envia request a oms
	// oms se encuentra en vm049
	addr = flag.String("addr", "10.6.46.59:50051", "ip address to connect to")

	//probar local
	//addr = flag.String("addr", "localhost:50051", "ip address to connect to")
	//condition = flag.String("condition", defaultcondition, "Condition to request")
)

func connectWithRetry() (*grpc.ClientConn, error) {
	for i := 0; i < 5; i++ { // retry up to 5 times
		// Create a connection to the server
		conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("Failed to connect: %v", err)
			if i == 4 { // if this was the fifth attempt, return the error
				return nil, err
			}
			log.Println("Retrying in 5 seconds...")
			time.Sleep(5 * time.Second) // wait for 5 seconds before retrying
		} else {
			return conn, nil // if connection succeeded, return the connection
		}
	}
	return nil, nil // this line should never be reached, but is required to satisfy the function signature
}

func main() {
	flag.Parse()

	/*
		    conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil{
				log.Fatalf("fallo la conexion: %v", err)
			}
			defer conn.Close()
	*/
	// fmt.Print("Solicite estado de 'Muerto' o 'Infectado': ")
	// scanner := bufio.NewScanner(os.Stdin)
	// scanner.Scan()
	// condition := strings.TrimSpace(scanner.Text())
	condition := "Muerto"
	// Verifica que la entrada sea válida
	if condition != "Muerto" && condition != "Infectado" {
		fmt.Println("Entrada no válida. Debe ingresar 'Muerto' o 'Infectado'.")
		return
	}

	conn, err := connectWithRetry()
	if err != nil {
		log.Fatalf("Failed to connect after 5 attempts: %v", err)
	}
	defer conn.Close()
	c := pb.NewRequestClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

    start := time.Now()
	r, err := c.RequestCondition(ctx, &pb.ConditionRequest{Condition: condition})
    end := time.Now()
    elapsed := end.Sub(start)
    fmt.Printf("tiempo de consutla: %s\n", elapsed)
	if err != nil {
		log.Fatalf("fallo en request: %v", err)
	}
	// Acceder a la lista de personas en la respuesta
	personas := r.GetPersons()

	// Iterar a través de la lista de personas y hacer algo con ellas
	for _, persona := range personas {
		nombre := persona.GetNombre()
		apellido := persona.GetApellido()
		// Hacer algo con el nombre y apellido, por ejemplo, imprimirlos
		fmt.Printf("Nombre: %s, Apellido: %s\n", nombre, apellido)
	}

}
