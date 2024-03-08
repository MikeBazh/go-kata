package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"
)

type SearchRequest struct {
	Query string `json:"query"`
}

type SearchResponse struct {
	Addresses []*Address `json:"addresses"`
}

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

func TestHandlerSearchByQuery(t *testing.T) {
	type args struct {
		link           string
		contentType    string
		wantMessage    []*Address
		wantStatusCode int
	}
	Address1 := Address{Value: "г Москва", UnrestrictedValue: "г Москва, ул Сухонская"}
	wantAddresses := []*Address{&Address1}
	tests := []struct {
		name string
		args args
	}{
		{"test1", args{"http://localhost:8080/api/address/search", "application/json", wantAddresses, http.StatusOK}},
		{"test2", args{"http://localhost:8080/api/address/search", "not json", wantAddresses, http.StatusBadRequest}},
		{"test3", args{"http://localhost:8080/api/address/123", "application/json", wantAddresses, 404}},
		// TODO: Add test cases.
	}
	go main()
	for _, tt := range tests {
		time.Sleep(time.Second / 5)
		// Определяем данные для нового поискового запроса
		newSearchRequest := SearchRequest{Query: "мск сухонск"}
		// Преобразуем структуру в JSON
		reqBody, err := json.Marshal(newSearchRequest)
		if err != nil {
			fmt.Println("Ошибка при кодировании JSON:", err)
			return
		}
		// Отправляем POST-запрос
		resp, err := http.Post(tt.args.link, tt.args.contentType, bytes.NewBuffer(reqBody))
		if err != nil {
			fmt.Println("Ошибка при выполнении POST-запроса:", err)
			return
		}
		fmt.Println(resp.StatusCode)
		defer resp.Body.Close()
		// Читаем ответ
		body, _ := io.ReadAll(resp.Body)
		newSearchResponse := SearchResponse{}
		err = json.Unmarshal(body, &newSearchResponse)
		if err != nil {
			if resp.StatusCode == http.StatusOK {
				fmt.Println("Ошибка получениия ответа:", err)
			}
			return
		}
		//fmt.Println("Получен ответ:", *newSearchResponse.Addresses[0])
		result := newSearchResponse.Addresses

		if resp.StatusCode != tt.args.wantStatusCode {
			t.Errorf("status want %d, got %d", tt.args.wantStatusCode, resp.StatusCode)
		}

		if !reflect.DeepEqual(tt.args.wantMessage[0].UnrestrictedValue, result[0].UnrestrictedValue) {
			t.Errorf("message want %v, got %v", tt.args.wantMessage[0].UnrestrictedValue, result[0].UnrestrictedValue)
		}
	}
}

type GeocodeRequest struct {
	Lat string
	Lng string
}

type GeocodeResponse struct {
	Addresses []*Address `json:"addresses"`
}

func TestHandlerSearchByGeo(t *testing.T) {
	type args struct {
		link           string
		contentType    string
		wantMessage    []*Address
		wantStatusCode int
	}
	Address1 := Address{Value: "г Москва", UnrestrictedValue: "108811, г Москва, поселение Московский, г Московский, Новомосковский округ, ул Бианки, д 4 к 2"}
	//Address2 := Address{Value: "г Москва", UnrestrictedValue: "г Московский, ул Бианки, д 4 к 2 108811"}
	//Address3 := Address{Value: "г Москва", UnrestrictedValue: "г Московский 108811"}
	wantAddresses := []*Address{&Address1}
	tests := []struct {
		name string
		args args
	}{
		{"test1", args{"http://localhost:8080/api/address/geocode", "application/json", wantAddresses, http.StatusOK}},
		{"test2", args{"http://localhost:8080/api/address/geocode", "not json", []*Address{}, http.StatusBadRequest}},
		{"test3", args{"http://localhost:8080/api/123/geocode", "application/json", []*Address{}, 404}},

		// TODO: Add test cases.
	}
	go main()
	for _, tt := range tests {
		time.Sleep(time.Second / 5)
		// Определяем данные для нового поискового запроса
		newSearchRequest := GeocodeRequest{
			Lat: "55.601983", Lng: "37.359486"}
		// Преобразуем структуру в JSON
		reqBody, err := json.Marshal(newSearchRequest)
		if err != nil {
			fmt.Println("Ошибка при кодировании JSON:", err)
			return
		}
		// Отправляем POST-запрос
		resp, err := http.Post(tt.args.link, tt.args.contentType, bytes.NewBuffer(reqBody))
		if err != nil {
			fmt.Println("Ошибка при выполнении POST-запроса:", err)
			return
		}
		fmt.Println(resp.StatusCode)
		defer resp.Body.Close()
		// Читаем ответ
		body, _ := io.ReadAll(resp.Body)
		newGeocodeResponse := GeocodeResponse{}
		err = json.Unmarshal(body, &newGeocodeResponse)
		if err != nil {
			if resp.StatusCode == http.StatusOK {
				fmt.Println("Ошибка получениия ответа:", err)
			}
			return
		}

		if resp.StatusCode != tt.args.wantStatusCode {
			t.Errorf("status want %d, got %d", tt.args.wantStatusCode, resp.StatusCode)
		}
		result := newGeocodeResponse.Addresses
		if !reflect.DeepEqual(tt.args.wantMessage[0].UnrestrictedValue, result[0].UnrestrictedValue) {
			t.Errorf("message want %v, got %v", tt.args.wantMessage[0].UnrestrictedValue, result[0].UnrestrictedValue)
		}
	}
}
