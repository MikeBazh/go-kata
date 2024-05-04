package main

import (
	"context"
	"fmt"
	"go-kata/2.server/5.server_http_api/geoservise_gRPC/internal/geoprovider/Dadata"
	"go-kata/2.server/5.server_http_api/geoservise_gRPC/internal/geoprovider/dto"
	pb "go-kata/2.server/5.server_http_api/geoservise_gRPC/internal/geoprovider/example"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) GeoCode(ctx context.Context, args *pb.GeoArgs) (*pb.GeoReply, error) {
	NewGeocodeResponse, err := Dadata.AskByGeo(args.Lat, args.Lon)
	if err != nil {
		return nil, err
	}
	return &pb.GeoReply{Addresses: ConvertAddresses(NewGeocodeResponse.Addresses)}, nil
}

func (s *server) AddressSearch(ctx context.Context, args *pb.QueryAddr) (*pb.SearchReply, error) {
	newSearchResponse, err := Dadata.AskByQuery(args.QueryAddr)
	if err != nil {
		return nil, err
	}
	return &pb.SearchReply{Addresses: ConvertAddresses(newSearchResponse.Addresses)}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	fmt.Println("server started on port:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func ConvertAddresses(addresses []*dto.Address) []*pb.Address {
	convertedAddresses := make([]*pb.Address, len(addresses))
	for i, addr := range addresses {
		convertedAddresses[i] = &pb.Address{
			Value:             addr.Value,
			UnrestrictedValue: addr.UnrestrictedValue,
			Data: &pb.Data{
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
