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

func main() {
	// Asignar variables si es que existe flag al compilar
	flag.Parse()

	// Crear connection por el mismo puerto del listener del servidor
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fallo la conexion: %v", err)
	}
	defer conn.Close()

	/*
		Utilizando los compilados de protocol buffer generar un cliente que
		pida el servicio Report definido en message.proto
		hacia el servidor con el cual se establecio la conexion conn.
	*/
	c := pb.NewReportClient(conn)

	/* Generar contexto
	Los contextos nos permiten compartir informacion entre distintos
	ambientes, en este caso el ambiente donde corre el cliente y el
	ambiente del servidor. El efecto de este codigo segun entiendo son
	los tiempos de ejecuccion presentes al ejecutarse tanto el cliente
	como el servidor. Ambos muestran su propio tiempo o context permite
	que ambos muestren tiempo del cliente? Sino, que informacion comparte
	este context?
	*/

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	/*
		Se realiza el request a traves de la conexion
		Notar que al compilar el .proto se crean structs de SeverityRequest y SeverityReply
		en message.pb.go, dicho archivo trata con las estructuras de datos y la serializacion, es decir,
		el ensamblado tangible de los datos para la comunicacion.

		Por otro lado message_grpc.pb.gp trata con la logica funcional de grpc, es decir, lo necesario
		para que cliente y servidor hablen el mismo idioma, en concreto, dar las herramientas que se
		llaman al utilizar el package pb que se llama explicitamente en el main.go tanto del cliente
		como del servidor.
	*/

	/*
		Al comprender los parrafos anteriores se entiende que se puede revisar el struct de SeverityRequest
		en el archivo adecuado, donde los campos definidos en el .proto son renombrados al mismo valor, pero con
		primera letra mayuscula.

		Al instanciar un struct, los cuales a veces funcionan como clases en go, los argumentos se pasan dentro de
		{} en vez de ()
	*/

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
