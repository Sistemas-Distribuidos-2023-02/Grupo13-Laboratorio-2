package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "github.com/VicenteRuizA/proto"
)

const (
	defaultname = "Cristiano Ronaldo"
	defaultcondition = "INFECTADO"
)

var (

    // NOTAR: se asume servidor en vm 50, de ahi la ip:puerto
    // ip se establece aqui, puerto en servidor.
	addr = flag.String("addr", "10.6.46.60:50051", "ip address to connect to")
	name = flag.String("name", defaultname, "Name to report")
	condition = flag.String("condition", defaultcondition, "Condition to report")
)


func main() {
	// Asignar variables si es que existe flag al compilar
	flag.Parse()

	
    // Crear connection por el mismo puerto del listener del servidor
    conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
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

    r, err := c.IdentifyCondition(ctx, &pb.SeverityRequest{Name : *name, Condition : *condition})
    if err != nil{
        log.Fatalf("could not greet: %v", err)
    }	
    log.Printf("Greeting: %s", r.GetMessage())
	
}
