package services

import (
	"context"
	"fmt"
	"go-kata/2.server/5.server_http_api/geoservise_gRPC/dto"
	example "go-kata/2.server/5.server_http_api/geoservise_gRPC/dto/example"
	pb "go-kata/2.server/5.server_http_api/geoservise_gRPC/dto/example"
	"google.golang.org/grpc"
)

type GeoproviderGrcp struct {
}

func (g *GeoproviderGrcp) AddressSearch(input string) ([]*dto.Address, error) {
	conn, err := grpc.Dial("geoprovider:50051", grpc.WithInsecure())
	var reply []*dto.Address
	if err != nil {
		fmt.Println("Ошибка при подключении к серверу:", err)
		return reply, err
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	args := &pb.QueryAddr{
		QueryAddr: input,
	}

	reply2, err := c.AddressSearch(context.Background(), args)
	if err != nil {
		fmt.Println("Ошибка при вызове удаленного метода:", err)
		return reply, err
	}
	reply = ConvertAddressesFromExample(reply2.Addresses)
	return reply, nil
}

func (g *GeoproviderGrcp) GeoCode(lat, lng string) ([]*dto.Address, error) {
	conn, err := grpc.Dial("geoprovider:50051", grpc.WithInsecure())
	var reply []*dto.Address
	if err != nil {
		fmt.Println("Ошибка при подключении к серверу:", err)
		return reply, err
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	args := &pb.GeoArgs{
		Lat: lat,
		Lon: lng,
	}

	reply2, err := c.GeoCode(context.Background(), args)
	if err != nil {
		fmt.Println("Ошибка при вызове удаленного метода:", err)
		return reply, err
	}
	reply = ConvertAddressesFromExample(reply2.Addresses)
	return reply, nil
}

func ConvertAddressesFromExample(addresses []*example.Address) []*dto.Address {
	convertedAddresses := make([]*dto.Address, len(addresses))
	for i, addr := range addresses {
		convertedAddresses[i] = &dto.Address{
			Value:             addr.Value,
			UnrestrictedValue: addr.UnrestrictedValue,
			Data: dto.Data{
				Source:               addr.Data.Source,
				Result:               addr.Data.Result,
				PostalCode:           addr.Data.PostalCode,
				Country:              addr.Data.Country,
				CountryIsoCode:       addr.Data.CountryIsoCode,
				FederalDistrict:      addr.Data.FederalDistrict,
				RegionIsoCode:        addr.Data.RegionIsoCode,
				RegionWithType:       addr.Data.RegionWithType,
				RegionType:           addr.Data.RegionType,
				RegionTypeFull:       addr.Data.RegionTypeFull,
				Region:               addr.Data.Region,
				AreaWithType:         addr.Data.AreaWithType,
				AreaType:             addr.Data.AreaType,
				AreaTypeFull:         addr.Data.AreaTypeFull,
				Area:                 addr.Data.Area,
				CityWithType:         addr.Data.CityWithType,
				CityType:             addr.Data.CityType,
				CityTypeFull:         addr.Data.CityTypeFull,
				City:                 addr.Data.City,
				CityArea:             addr.Data.CityArea,
				CityDistrictWithType: addr.Data.CityDistrictWithType,
				CityDistrictType:     addr.Data.CityDistrictType,
				CityDistrictTypeFull: addr.Data.CityDistrictTypeFull,
				CityDistrict:         addr.Data.CityDistrict,
				SettlementWithType:   addr.Data.SettlementWithType,
				SettlementType:       addr.Data.SettlementType,
				SettlementTypeFull:   addr.Data.SettlementTypeFull,
				Settlement:           addr.Data.Settlement,
				StreetWithType:       addr.Data.StreetWithType,
				StreetType:           addr.Data.StreetType,
				StreetTypeFull:       addr.Data.StreetTypeFull,
				Street:               addr.Data.Street,
				HouseType:            addr.Data.HouseType,
				HouseTypeFull:        addr.Data.HouseTypeFull,
				House:                addr.Data.House,
				BlockType:            addr.Data.BlockType,
				BlockTypeFull:        addr.Data.BlockTypeFull,
				Block:                addr.Data.Block,
				Entrance:             addr.Data.Entrance,
				Floor:                addr.Data.Floor,
				FlatFiasId:           addr.Data.FlatFiasId,
				FlatType:             addr.Data.FlatType,
				FlatTypeFull:         addr.Data.FlatTypeFull,
				Flat:                 addr.Data.Flat,
				FlatArea:             addr.Data.FlatArea,
				SquareMeterPrice:     addr.Data.SquareMeterPrice,
				FlatPrice:            addr.Data.FlatPrice,
				PostalBox:            addr.Data.PostalBox,
				FiasCode:             addr.Data.FiasCode,
				FiasLevel:            addr.Data.FiasLevel,
				FiasActualityState:   addr.Data.FiasActualityState,
				CapitalMarker:        addr.Data.CapitalMarker,
				Okato:                addr.Data.Okato,
				Oktmo:                addr.Data.Oktmo,
				TaxOffice:            addr.Data.TaxOffice,
				TaxOfficeLegal:       addr.Data.TaxOfficeLegal,
				Timezone:             addr.Data.Timezone,
				GeoLat:               addr.Data.GeoLat,
				GeoLon:               addr.Data.GeoLon,
				BeltwayHit:           addr.Data.BeltwayHit,
				BeltwayDistance:      addr.Data.BeltwayDistance,
				UnparsedParts:        addr.Data.UnparsedParts,
			},
		}
	}
	return convertedAddresses
}
