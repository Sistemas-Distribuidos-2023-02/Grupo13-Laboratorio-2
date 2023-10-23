package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	pb "github.com/VicenteRuizA/proto_lab2/oms_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// const (
// 	defaultname = "Cristiano"
//     defaultsurname = "Ronaldo"
// 	defaultcondition = "INFECTADO"
// )

var (

	// latinoamerica envia mensajes a oms
	// oms se encuentra en vm049
	addr = flag.String("addr", "10.6.46.59:50051", "ip address to connect to")

	//probar local
	//addr = flag.String("addr", "localhost:50051", "ip address to connect to")
	// name      = flag.String("name", defaultname, "Name to report")
	// surname   = flag.String("surname", defaultsurname, "Surname to report")
	// condition = flag.String("condition", defaultcondition, "Condition to report")
)

// lee los nombres del archivos names.txt
func readNamesFromFile(filename string, numNames int) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var names []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		names = append(names, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(names), func(i, j int) {
		names[i], names[j] = names[j], names[i]
	})

	return names[:numNames]
}

func separarNombreApellido(nombreCompleto string) (nombre, apellido string, err error) {
	// Divide el string usando el espacio como delimitador
	parts := strings.Split(nombreCompleto, " ")

	// Verifica si hay al menos dos partes (nombre y apellido)
	if len(parts) < 2 {
		err = fmt.Errorf("El formato del nombre completo no es válido")
		return
	}

	// Asigna el nombre y apellido
	nombre = parts[0]
	apellido = strings.Join(parts[1:], " ")

	return
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
	// Asignar variables si es que existe flag al compilar
	flag.Parse()

	// Crear connection por el mismo puerto del listener del servidor
	/*
		conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("fallo la conexion: %v", err)
		}
		defer conn.Close()
	*/

	conn, err := connectWithRetry()
	if err != nil {
		log.Fatalf("Failed to connect after 5 attempts: %v", err)
	}
	defer conn.Close()

	c := pb.NewReportClient(conn)

	//aquí se amplio tiempo
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	names := readNamesFromFile("names.txt", 5) // Obtén 5 nombres al azar al inicio
	// nombre, apellido, err := separarNombreApellido(names[0])
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	for _, name := range names {
		state := "Infectado"
		if rand.Float32() > 0.55 {
			state = "Muerto"
		}

		log.Printf("Nombre: %v", name)
		log.Printf("Estado: %v", state)

		nombre, apellido, err := separarNombreApellido(name)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		//_, err := client.UploadUser(context.Background(), user)
		//Cristiano RonaldoInfectado
		r, err := c.IdentifyCondition(ctx, &pb.SeverityRequest{Name: nombre, Surname: apellido, Condition: state})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.GetMessage())
	}
	for {

		name := readNamesFromFile("names.txt", 1)
		state := "Infectado"
		if rand.Float32() > 0.55 {
			state = "Muerto"
		}
		log.Printf("Selected name: %s, Estado: %s", name[0], state)

		// Realiza alguna operación con el nombre, por ejemplo, enviarlo al servidor gRPC
		nombre, apellido, err := separarNombreApellido(name[0])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		r, err := c.IdentifyCondition(ctx, &pb.SeverityRequest{Name: nombre, Surname: apellido, Condition: state})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.GetMessage())
		time.Sleep(3 * time.Second) // Espera 3 segundos antes de seleccionar el próximo nombre
	}
	// state := "Infectado"
	// if rand.Float32() > 0.55 {
	// 	state = "Muerto"
	// }
	// log.Printf("Estado: %v", state)

	// r, err := c.IdentifyCondition(ctx, &pb.SeverityRequest{Name: nombre, Surname: apellido, Condition: state})
	// if err != nil {
	// 	log.Fatalf("could not greet: %v", err)
	// }
	// log.Printf("Greeting: %s", r.GetMessage())

}
