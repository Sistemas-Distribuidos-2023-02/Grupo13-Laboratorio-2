package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
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

func getUserInput() string {
	fmt.Print("Solicite estado de 'Muerto' o 'Infectado', o escriba 'salir' para salir: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

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

	conn, err := connectWithRetry()
	if err != nil {
		log.Fatalf("Failed to connect after 5 attempts: %v", err)
	}
	defer conn.Close()
	c := pb.NewRequestClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	for {
		condition := getUserInput()

		// Verificar si el usuario quiere salir del programa
		if strings.ToLower(condition) == "salir" {
			fmt.Println("Saliendo del programa...")
			break // Salir del bucle y terminar el programa
		}

		// Verificar la entrada del usuario y continuar con el código principal
		if condition != "Muerto" && condition != "Infectado" {
			fmt.Println("Entrada no válida. Debe ingresar 'Muerto' o 'Infectado'.")
			continue // Volver a solicitar la entrada del usuario
		}

		r, err := c.RequestCondition(ctx, &pb.ConditionRequest{Condition: condition})
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
}
