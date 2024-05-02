package main

import (
	"go-kata/2.server/5.server_http_api/geoservise_rpc/internal/geoprovider/Dadata"
	"go-kata/2.server/5.server_http_api/geoservise_rpc/internal/geoprovider/dto"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type GeoProvider interface {
	AddressSearch(input string) ([]*dto.Address, error)
	GeoCode(lat, lng string) ([]*dto.Address, error)
}

// GeoService реализует интерфейс GeoProvider и предоставляет удаленные процедуры для работы с провайдерами геокодирования.
type GeoService struct {
	// Здесь могут быть поля для хранения информации о провайдерах и других конфигурационных данных.
}

// AddressSearch выполняет поиск адреса с помощью указанного провайдера геокодирования.
func (g *GeoService) AddressSearch(args string, reply *[]*dto.Address) error {
	newSearchResponse, err := Dadata.AskByQuery(args)
	if err != nil {
		return err
	}
	//fmt.Println("newSearchResponse:", newSearchResponse)
	*reply = newSearchResponse.Addresses
	//fmt.Println(reply)
	return nil
}

type GeoArgs struct {
	Lat string
	Lon string
}

// GeoCode выполняет геокодирование с помощью указанного провайдера геокодирования.
func (g *GeoService) GeoCode(args *dto.GeoArgs, reply *[]*dto.Address) error {
	newSearchResponse, err := Dadata.AskByGeo(args.Lat, args.Lon)
	if err != nil {
		return err
	}
	//fmt.Println("newSearchResponse:", newSearchResponse)
	*reply = newSearchResponse.Addresses
	//fmt.Println(reply)
	return nil
}

func main() {
	geoService := new(GeoService)
	rpc.Register(geoService)

	listener, err := net.Listen("tcp", ":8070")
	if err != nil {
		log.Fatal("Error starting RPC server:", err)
	}
	log.Println("RPC server started on port 8070")

	listenerJsonRpc, err := net.Listen("tcp", ":8060")
	if err != nil {
		log.Fatal("Listen error:", err)
	}
	log.Println("JsonRPC server started on port 8060")

	go func() {
		for {
			// Ожидание соединений на любом из портов
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Accept error:", err)
				continue
			}
			go func(conn net.Conn) {
				defer conn.Close()
				rpc.ServeConn(conn)
			}(conn)
		}
	}()

	go func() {
		for {
			connJsonRpc, err := listenerJsonRpc.Accept()
			if err != nil {
				log.Println("Accept error:", err)
				continue
			}

			go func(connJsonRpc net.Conn) {
				defer connJsonRpc.Close()
				jsonrpc.ServeConn(connJsonRpc)
			}(connJsonRpc)
		}
	}()

	select {} // для предотвращения завершения программы

}

//	geoService := new(GeoService)
//	rpc.Register(geoService)
//
//	listenerRPC, err := net.Listen("tcp", ":8070")
//	if err != nil {
//		log.Fatal("Error starting RPC server:", err)
//	}
//	defer listenerRPC.Close()
//	log.Println("RPC server started on port 8070")
//
//	listenerJSONRPC, err := net.Listen("tcp", ":8060")
//	if err != nil {
//		log.Fatal("Error starting JSON-RPC server:", err)
//	}
//	defer listenerJSONRPC.Close()
//	log.Println("JSON-RPC server started on port 8060")
//
//	go func() {
//		for {
//			conn, err := listenerRPC.Accept()
//			if err != nil {
//				log.Println("Accept error:", err)
//				continue
//			}
//
//			go func(conn net.Conn) {
//				defer conn.Close()
//				rpc.ServeConn(conn)
//			}(conn)
//		}
//	}()
//
//	go func() {
//		for {
//			conn, err := listenerJSONRPC.Accept()
//			if err != nil {
//				log.Println("Accept error:", err)
//				continue
//			}
//
//			go func(conn net.Conn) {
//				defer conn.Close()
//				jsonrpc.ServeConn(conn)
//			}(conn)
//		}
//	}()
//
//	select {} // Бесконечный цикл для предотвращения завершения программы
//}

//err := godotenv.Load()
//if err != nil {
//log.Fatal("Error load .env file: ", err)
//}
//protocol := os.Getenv("RPC_PROTOCOL")
//
////if !(protocol == "json-rpc" || protocol == "rpc") {
////	log.Fatalf("Server error: Protocol must be set to 'json-rpc' or 'rpc'")
////}
//
//geoService := new(GeoService)
//rpc.Register(geoService)
//
//if protocol == "rpc" {
//listener, err := net.Listen("tcp", ":8070")
//if err != nil {
//log.Fatal("Error starting RPC server:", err)
//}
//log.Println("RPC server started on port 8070")
//for {
//conn, err := listener.Accept()
//if err != nil {
//log.Fatal("Accept error:", err)
//}
//
//go func(conn net.Conn) {
//defer conn.Close()
//rpc.ServeConn(conn)
//}(conn)
//}
//}
//
//if protocol == "json-rpc" {
//listenerJsonRpc, err := net.Listen("tcp", ":8060")
//if err != nil {
//log.Fatal("Listen error:", err)
//}
//log.Println("JsonRPC server started on port 8060")
//for {
//connJsonRpc, err := listenerJsonRpc.Accept()
//if err != nil {
//log.Fatal("Accept error:", err)
//}
//
//go func(connJsonRpc net.Conn) {
//defer connJsonRpc.Close()
//jsonrpc.ServeConn(connJsonRpc)
//}(connJsonRpc)
//}
//}
//}
