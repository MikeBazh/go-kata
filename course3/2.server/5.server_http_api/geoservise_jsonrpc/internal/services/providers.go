package services

import (
	"fmt"
	"go-kata/2.server/5.server_http_api/geoservise_jsonrpc/dto"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type GeoProvider interface {
	AddressSearch(input string) ([]*dto.Address, error)
	GeoCode(lat, lng string) ([]*dto.Address, error)
}

type GeoProviderRcp struct {
}

func NewGeoProvider(protocol string) GeoProvider {
	switch protocol {
	case "json-rpc":
		return &GeoProviderJsonRcp{}
	case "rpc":
		return &GeoProviderRcp{}
	default:
		return nil
	}
}

func (g *GeoProviderRcp) AddressSearch(input string) ([]*dto.Address, error) {
	client, err := rpc.Dial("tcp", "geoprovider:8070")
	var reply []*dto.Address
	if err != nil {
		fmt.Println("Ошибка при подключении к серверу:", err)
		return reply, err
	}
	err = client.Call("GeoService.AddressSearch", input, &reply)
	if err != nil {
		fmt.Println("Ошибка при вызове удаленного метода:", err)
		return reply, err
	}
	return reply, nil
}

func (g *GeoProviderRcp) GeoCode(lat, lng string) ([]*dto.Address, error) {

	var reply []*dto.Address
	client, err := rpc.Dial("tcp", "geoprovider:8070")
	if err != nil {
		fmt.Println("Ошибка при подключении к серверу:", err)
		return reply, err
	}

	args := dto.GeoArgs{lat, lng}
	err = client.Call("GeoService.GeoCode", args, &reply)
	if err != nil {
		fmt.Println("Ошибка при вызове удаленного метода:", err)
		return reply, err
	}
	return reply, nil
}

type GeoProviderJsonRcp struct {
}

func (g *GeoProviderJsonRcp) AddressSearch(input string) ([]*dto.Address, error) {
	client, err := jsonrpc.Dial("tcp", "geoprovider:8060")
	var reply []*dto.Address
	if err != nil {
		fmt.Println("Ошибка при подключении к серверу:", err)
		return reply, err
	}
	err = client.Call("GeoService.AddressSearch", input, &reply)
	if err != nil {
		fmt.Println("Ошибка при вызове удаленного метода:", err)
		return reply, err
	}
	return reply, nil
}

func (g *GeoProviderJsonRcp) GeoCode(lat, lng string) ([]*dto.Address, error) {

	var reply []*dto.Address
	client, err := jsonrpc.Dial("tcp", "geoprovider:8060")
	if err != nil {
		fmt.Println("Ошибка при подключении к серверу:", err)
		return reply, err
	}

	args := dto.GeoArgs{lat, lng}
	err = client.Call("GeoService.GeoCode", args, &reply)
	if err != nil {
		fmt.Println("Ошибка при вызове удаленного метода:", err)
		return reply, err
	}
	return reply, nil
}
