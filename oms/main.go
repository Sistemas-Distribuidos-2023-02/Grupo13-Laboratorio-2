package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
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

type DataRecord struct {
	ID       string
	Datanote string
	Estado   string
}

func LeerArchivoData(filename string) ([]DataRecord, error) {
	var records []DataRecord

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "-")
		if len(fields) == 3 {
			id := fields[0]
			datanote := fields[1]
			estado := fields[2]

			record := DataRecord{
				ID:       id,
				Datanote: datanote,
				Estado:   estado,
			}

			records = append(records, record)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

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

func ListadeID(condition string) ([]string, []string) {
	records, err := LeerArchivoData("DATA.txt")
	if err != nil {
		log.Fatal(err)
	}

	var idData1 []string
	var idData2 []string

	for _, record := range records {
		fmt.Printf("ID: %s, Datanote: %s, Estado: %s\n", record.ID, record.Datanote, record.Estado)

		if record.Estado == condition {
			switch record.Datanote {
			case "datanode1":
				idData1 = append(idData1, record.ID)
			case "datanode2":
				idData2 = append(idData2, record.ID)
			}
		}
	}
	return idData1, idData2
}

func connectWithRetry(addr string) (*grpc.ClientConn, error) {
	for i := 0; i < 5; i++ { // retry up to 5 times
		// Create a connection to the server
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

type server struct {
	pbs.UnimplementedReportServer
	pbs.UnimplementedRequestServer
}

func (s *server) IdentifyCondition(ctx context.Context, in *pbs.SeverityRequest) (*pbs.SeverityReply, error) {
	log.Printf("Received: \n Nombre: %v\n Apellido: %v\n Condicion de %v", in.GetName(), in.GetSurname(), in.GetCondition())

	addr := conexionADatanode(in.GetSurname(), in.GetCondition(), name_id)

	/*
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("fallo la conexion: %v", err)
		}
		defer conn.Close()
	*/
	conn, err := connectWithRetry(addr)
	if err != nil {
		log.Fatalf("Failed to connect after 5 attempts: %v", err)
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

// func (s *server) RequestCondition(ctx context.Context, in *pbs.ConditionRequest) (*pbs.ConditionReply, error) {
// 	log.Printf("Received: Condicion %v", in.GetCondition())

// 	// Cambiar por conexion segun datanode que contenga ids de condicion solicitada
// 	//addr := conexionADatanode(name, in.GetCondition(), name_id)

// 	/*
// 		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 		if err != nil {
// 			log.Fatalf("fallo la conexion: %v", err)
// 		}
// 		defer conn.Close()
// 	*/
// 	addr1 := "10.6.46.61:50051"
// 	addr2 := "10.6.46.62:50051"
// 	listaData1, listaData2 := ListadeID(in.GetCondition())
// 	//var listaONU []string
// 	var listaONU []*pbs.Person

// 	//conexion a dn1
// 	conn_dn1, err := connectWithRetry(addr1)
// 	if err != nil {
// 		log.Fatalf("Failed to connect after 5 attempts: %v", err)
// 	}
// 	defer conn_dn1.Close()
// 	c := pbc.NewLoadClient(conn_dn1)

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	//conexion a dn2
// 	conn_dn2, err := connectWithRetry(addr2)
// 	if err != nil {
// 		log.Fatalf("Failed to connect after 5 attempts: %v", err)
// 	}
// 	defer conn_dn2.Close()
// 	c2 := pbc.NewLoadClient(conn_dn2)

// 	ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	for _, idea := range listaData1 {
// 		r, err := c.RequestData(ctx, &pbc.DataRequest{Id: idea})
// 		if err != nil {
// 			log.Fatalf("could not greet: %v", err)
// 		}
// 		log.Printf("Reply from server: nombre: %s apellido: %s", r.Nombre, r.Apellido)

//         person := &pbs.Person{
// 		    Nombre:   r.Nombre,
// 		    Apellido: r.Apellido,
// 	    }

// 		listaONU = append(listaONU, person)
// 	}

// 	for _, idea := range listaData2 {
// 		r, err := c2.RequestData(ctx2, &pbc.DataRequest{Id: idea})
// 		if err != nil {
// 			log.Fatalf("could not greet: %v", err)
// 		}
// 		log.Printf("Reply from server: nombre: %s apellido: %s", r.Nombre, r.Apellido)
// 		//lista a mandar a ONU

//         person := &pbs.Person{
// 		    Nombre:   r.Nombre,
// 		    Apellido: r.Apellido,
// 	    }

// 		listaONU = append(listaONU, person)

// 	}
// 	// conn, err := connectWithRetry(addr)
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to connect after 5 attempts: %v", err)
// 	// }
// 	// defer conn.Close()
// 	// c := pbc.NewLoadClient(conn)

// 	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	// defer cancel()

// 	// r, err := c.RequestData(ctx, &pbc.DataRequest{Id: strconv.Itoa(name_id)})
// 	// if err != nil {
// 	// 	log.Fatalf("could not greet: %v", err)
// 	// }
// 	// log.Printf("Reply from server: nombre: %s apellido: %s", r.Nombre, r.Apellido)

// 	return &pbs.ConditionReply{Persons: listaONU}, nil
// }

func startServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pbs.RegisterReportServer(s, &server{})
	pbs.RegisterRequestServer(s, &server{})
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
