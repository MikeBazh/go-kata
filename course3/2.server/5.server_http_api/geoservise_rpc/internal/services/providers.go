package services

import (
	"fmt"
	"go-kata/2.server/5.server_http_api/geoservise_rpc/dto"
	"net/rpc"
)

type GeoProvider interface {
	AddressSearch(input string) ([]*dto.Address, error)
	GeoCode(lat, lng string) ([]*dto.Address, error)
}

type GeoProviderDadata struct {
}

func NewGeoProviderDadata() GeoProvider {
	return &GeoProviderDadata{}
}

func (g *GeoProviderDadata) AddressSearch(input string) ([]*dto.Address, error) {
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

func (g *GeoProviderDadata) GeoCode(lat, lng string) ([]*dto.Address, error) {

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
