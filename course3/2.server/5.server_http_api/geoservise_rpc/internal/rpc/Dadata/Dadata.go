package Dadata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ekomobile/dadata/v2"
	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/ekomobile/dadata/v2/client"
	"go-kata/2.server/5.server_http_api/geoservise_rpc/internal/rpc/dto"
	"io"
	"net/http"
)

func AskByQuery(queryAddr string) (dto.SearchResponse, error) {
	var Addresses []*dto.Address
	creds := client.Credentials{
		ApiKeyValue:    "f1487054fd3354afa7fff0a8aed0276d6b9b1545",
		SecretKeyValue: "91662072e8c12d1185984451eece3124af4cb800",
	}
	api := dadata.NewSuggestApi(client.WithCredentialProvider(&creds))
	params := suggest.RequestParams{
		Query: queryAddr,
	}
	suggestions, err := api.Address(context.Background(), &params)
	if err != nil {
		return dto.SearchResponse{}, err
	}
	//newSearchResponse := dto.SearchResponse{}
	//fmt.Println(suggestions)
	for _, s := range suggestions {
		//fmt.Println(s.Value)
		newData := dto.Data{Source: s.Data.Source, Result: s.Data.Result, PostalCode: s.Data.PostalCode, Country: s.Data.Country, CountryIsoCode: s.Data.CountryIsoCode, FederalDistrict: s.Data.FederalDistrict, Region: s.Data.Region, City: s.Data.City, Street: s.Data.Street}
		Addresses = append(Addresses, &dto.Address{Value: s.Value, UnrestrictedValue: s.UnrestrictedValue, Data: newData})
	}

	var Response dto.SearchResponse
	Response.Addresses = Addresses
	return Response, nil
}

func AskByGeo(lat, lon string) (s dto.GeocodeResponse, err error) {
	// Укажем URL и тело запроса
	url := "http://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address"
	body := map[string]string{"lat": lat, "lon": lon, "radius_meters": "100"}

	// Кодируем тело запроса в JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Ошибка AskDadata при кодировании JSON:", err)
		return dto.GeocodeResponse{}, err
	}

	NewClient := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Ошибка AskDadata при создании запроса:", err)
		return dto.GeocodeResponse{}, err
	}

	// Устанавливаем заголовки запроса
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Token f1487054fd3354afa7fff0a8aed0276d6b9b1545")

	resp, err := NewClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка AskDadata при выполнении запроса:", err)
		return dto.GeocodeResponse{}, err
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	respBody, err := io.ReadAll(resp.Body)
	NewGeocodeResponse := dto.GeocodeResponseSuggest{}
	err = json.Unmarshal(respBody, &NewGeocodeResponse)
	if err != nil {
		fmt.Println("Ошибка AskDadata при чтении тела ответа:", err)
		return dto.GeocodeResponse{}, err
	}
	var Response dto.GeocodeResponse
	Response.Addresses = NewGeocodeResponse.Addresses
	return Response, nil
}

//type GeocodeResponse struct {
//	Addresses []*Address `json:"addresses"`
//}
//
//type SearchResponse struct {
//	Addresses []*Address `json:"addresses"`
//}
//
//type GeocodeResponseSuggest struct {
//	Addresses []*Address `json:"suggestions"`
//}
//
//type SearchResponseSuggest struct {
//	Addresses []*Address `json:"suggestions"`
//}

type Address struct {
	Value             string `json:"value"`
	UnrestrictedValue string `json:"unrestricted_value"`
	Data              Data   `json:"data"`
}

type Data struct {
	Source                 string      `json:"source"`
	Result                 string      `json:"result"`
	PostalCode             string      `json:"postal_code"`
	Country                string      `json:"country"`
	CountryIsoCode         string      `json:"country_iso_code"`
	FederalDistrict        string      `json:"federal_district"`
	RegionFiasID           string      `json:"region_fias_id"`
	RegionKladrID          string      `json:"region_kladr_id"`
	RegionIsoCode          string      `json:"region_iso_code"`
	RegionWithType         string      `json:"region_with_type"`
	RegionType             string      `json:"region_type"`
	RegionTypeFull         string      `json:"region_type_full"`
	Region                 string      `json:"region"`
	AreaFiasID             string      `json:"area_fias_id"`
	AreaKladrID            string      `json:"area_kladr_id"`
	AreaWithType           string      `json:"area_with_type"`
	AreaType               string      `json:"area_type"`
	AreaTypeFull           string      `json:"area_type_full"`
	Area                   string      `json:"area"`
	CityFiasID             string      `json:"city_fias_id"`
	CityKladrID            string      `json:"city_kladr_id"`
	CityWithType           string      `json:"city_with_type"`
	CityType               string      `json:"city_type"`
	CityTypeFull           string      `json:"city_type_full"`
	City                   string      `json:"city"`
	CityArea               string      `json:"city_area"`
	CityDistrictFiasID     string      `json:"city_district_fias_id"`
	CityDistrictKladrID    string      `json:"city_district_kladr_id"`
	CityDistrictWithType   string      `json:"city_district_with_type"`
	CityDistrictType       string      `json:"city_district_type"`
	CityDistrictTypeFull   string      `json:"city_district_type_full"`
	CityDistrict           string      `json:"city_district"`
	SettlementFiasID       string      `json:"settlement_fias_id"`
	SettlementKladrID      string      `json:"settlement_kladr_id"`
	SettlementWithType     string      `json:"settlement_with_type"`
	SettlementType         string      `json:"settlement_type"`
	SettlementTypeFull     string      `json:"settlement_type_full"`
	Settlement             string      `json:"settlement"`
	StreetFiasID           string      `json:"street_fias_id"`
	StreetKladrID          string      `json:"street_kladr_id"`
	StreetWithType         string      `json:"street_with_type"`
	StreetType             string      `json:"street_type"`
	StreetTypeFull         string      `json:"street_type_full"`
	Street                 string      `json:"street"`
	HouseFiasID            string      `json:"house_fias_id"`
	HouseKladrID           string      `json:"house_kladr_id"`
	HouseType              string      `json:"house_type"`
	HouseTypeFull          string      `json:"house_type_full"`
	House                  string      `json:"house"`
	HouseCadNum            string      `json:"house_cadnum"`
	BlockType              string      `json:"block_type"`
	BlockTypeFull          string      `json:"block_type_full"`
	Block                  string      `json:"block"`
	Entrance               string      `json:"entrance"`
	Floor                  string      `json:"floor"`
	FlatFiasId             string      `json:"flat_fias_id"`
	FlatType               string      `json:"flat_type"`
	FlatTypeFull           string      `json:"flat_type_full"`
	Flat                   string      `json:"flat"`
	FlatArea               string      `json:"flat_area"`
	FlatCadNum             string      `json:"flat_cadnum"`
	SquareMeterPrice       string      `json:"square_meter_price"`
	FlatPrice              string      `json:"flat_price"`
	PostalBox              string      `json:"postal_box"`
	FiasID                 string      `json:"fias_id"`
	FiasCode               string      `json:"fias_code"`
	FiasLevel              string      `json:"fias_level"`
	FiasActualityState     string      `json:"fias_actuality_state"`
	KladrID                string      `json:"kladr_id"`
	CapitalMarker          string      `json:"capital_marker"`
	Okato                  string      `json:"okato"`
	Oktmo                  string      `json:"oktmo"`
	TaxOffice              string      `json:"tax_office"`
	TaxOfficeLegal         string      `json:"tax_office_legal"`
	Timezone               string      `json:"timezone"`
	GeoLat                 string      `json:"geo_lat"`
	GeoLon                 string      `json:"geo_lon"`
	BeltwayHit             string      `json:"beltway_hit"`
	BeltwayDistance        string      `json:"beltway_distance"`
	QualityCodeGeoRaw      interface{} `json:"qc_geo"`
	QualityCodeCompleteRaw interface{} `json:"qc_complete"`
	QualityCodeHouseRaw    interface{} `json:"qc_house"`
	QualityCodeRaw         interface{} `json:"qc"`
	UnparsedParts          string      `json:"unparsed_parts"`
	Metro                  struct{}    `json:"metro"`
}
