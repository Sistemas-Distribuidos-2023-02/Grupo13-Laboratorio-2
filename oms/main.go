package main

import (
	"context"
	"flag"
    "fmt"
	"log"
    "net"
    "time"
    "sync"
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

	// Use fmt.Sprintf to format the string with variables.
	replyMessage := fmt.Sprintf("Se ha reportado exitosamente que %s %s esta %s", in.GetName(), in.GetSurname(), in.GetCondition())

	return &pbs.SeverityReply{Message: replyMessage}, nil
}

func startServer(wg *sync.WaitGroup) {
    defer wg.Done()
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
// client -------------------------------------------------------------------------------------------

const (
	defaultname = "Cristiano"
    defaultsurname = "Ronaldo"
	defaultcondition = "INFECTADO"
)

var (
    //probar en en vms donde oms, servidor al cual pide requiest, esta en 049
    //addr = flag.String("addr", "10.6.46.59:50051", "ip address to connect to")

    //probar local
    addr = flag.String("addr", "localhost:50051", "ip address to connect to")
	name = flag.String("name", defaultname, "Name to report")
	surname = flag.String("surname", defaultsurname, "Surname to report")
	condition = flag.String("condition", defaultcondition, "Condition to report")
)

func startClient(wg *sync.WaitGroup){
    defer wg.Done()
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
	c := pbc.NewSaveClient(conn)

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

    r, err := c.SaveNaming(ctx, &pbc.SaveRequest{Id : *id, Name : *name, Surname : *surname})
    if err != nil{
        log.Fatalf("could not greet: %v", err)
    }	
    log.Printf("Reply from server: %s", r.GetMessage())
	
}
// main -------------------------------------------------------------------------------------------
func main() {
	flag.Parse()
    var wg sync.WaitGroup
    wg.Add(2) // We are starting two goroutines
    go startClient(&wg)
    go startServer(&wg)
    wg.Wait() // This will block until both goroutines have called Done
}
