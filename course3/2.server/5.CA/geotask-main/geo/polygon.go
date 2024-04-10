package geo

import (
	geo "github.com/kellydunn/golang-geo"
	"math/rand"
	"time"
)

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type PolygonChecker interface {
	Contains(point Point) bool // проверить, находится ли точка внутри полигона
	Allowed() bool             // разрешено ли входить в полигон
	RandomPoint() Point        // сгенерировать случайную точку внутри полигона
}

type Polygon struct {
	polygon *geo.Polygon
	allowed bool
}

func (p Polygon) Contains(point Point) bool {
	geoPoint := geo.NewPoint(point.Lat, point.Lng)
	return p.polygon.Contains(geoPoint)
}

func (p Polygon) Allowed() bool {
	return p.allowed
}

func (p Polygon) RandomPoint() Point {
	// Generate random points until a point inside the polygon is found
	//log.Println("Lat max, min:", p.MaxLat(), p.MinLat())
	//log.Println("Lng max, min:", p.MaxLng(), p.MinLng())
	randSource := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(randSource)
	for {
		// Generate random latitude and longitude within the bounding box of the polygon
		randomLat := rnd.Float64()*(p.MaxLat()-p.MinLat()) + p.MinLat()
		randomLng := rnd.Float64()*(p.MaxLng()-p.MinLng()) + p.MinLng()
		//log.Println(p.Contains(Point{Lat: randomLat, Lng: randomLng}))
		// Check if the random point is inside the polygon
		if p.Contains(Point{Lat: randomLat, Lng: randomLng}) {
			return Point{Lat: randomLat, Lng: randomLng}
		}
	}
}

func NewPolygon(points []Point, allowed bool) *Polygon {
	// используем библиотеку golang-geo для создания полигона
	var geoPoints []*geo.Point
	//var p geo.Point
	for _, v := range points {

		geoPoints = append(geoPoints, geo.NewPoint(v.Lat, v.Lng))
	}
	return &Polygon{
		polygon: geo.NewPolygon(geoPoints),
		allowed: allowed,
	}
}

func CheckPointIsAllowed(point Point, allowedZone PolygonChecker, disabledZones []PolygonChecker) bool {
	// проверить, находится ли точка в разрешенной зоне
	var checked = true
	for _, disabledZone := range disabledZones {
		if disabledZone.Contains(point) {
			checked = false
		}
	}
	if !allowedZone.Contains(point) {
		checked = false
	}
	return checked
}

func GetRandomAllowedLocation(allowedZone PolygonChecker, disabledZones []PolygonChecker) Point {
	var point Point
	// получение случайной точки в разрешенной зоне
	for {
		point = allowedZone.RandomPoint()
		if CheckPointIsAllowed(point, allowedZone, disabledZones) {
			return point
		}
	}
}

// полигоны лежат в /public/js/polygons.js
func NewDisAllowedZone1() *Polygon {
	points := []Point{
		{Lat: 59.902742187627325, Lng: 30.35368172093575},
		{Lat: 59.90015959974209, Lng: 30.41290489598458},
		{Lat: 59.842429456164574, Lng: 30.411531604968953},
		{Lat: 59.836047143247896, Lng: 30.373766102039266},
	}
	return NewPolygon(points, false)
}

func NewDisAllowedZone2() *Polygon {
	points := []Point{
		{Lat: 60.051063834232714, Lng: 30.28244720269174},
		{Lat: 60.0509781359604, Lng: 30.341498716363613},
		{Lat: 60.02036963316746, Lng: 30.363471372613613},
		{Lat: 60.01650940538451, Lng: 30.31986938286752},
	}

	return NewPolygon(points, false)
}

func NewAllowedZone() *Polygon {
	// Define the vertices of the polygon
	points := []Point{
		{Lat: 60.05759504176843, Lng: 30.14495968779295},
		{Lat: 60.07986463778022, Lng: 30.190278291308577},
		{Lat: 60.08269008837324, Lng: 30.20143628081053},
		{Lat: 60.08410272287511, Lng: 30.21662831267088},
		{Lat: 60.08620015941349, Lng: 30.229760408007795},
		{Lat: 60.09210650847325, Lng: 30.245724916064436},
		{Lat: 60.09480253335778, Lng: 30.252677201831037},
		{Lat: 60.09681370986395, Lng: 30.26031613310545},
		{Lat: 60.09883394571154, Lng: 30.272926706427032},
		{Lat: 60.09897148869136, Lng: 30.28416398873299},
		{Lat: 60.0954811436399, Lng: 30.3286112095949},
		{Lat: 60.09327428225354, Lng: 30.363445393347437},
		{Lat: 60.086444857223825, Lng: 30.376478582226927},
		{Lat: 60.064253238880035, Lng: 30.385181009375746},
		{Lat: 60.055490095341256, Lng: 30.39468944033354},
		{Lat: 60.04344323015362, Lng: 30.437052249514753},
		{Lat: 60.03374429411284, Lng: 30.44212698897093},
		{Lat: 60.01845570695627, Lng: 30.45914292296141},
		{Lat: 60.009144281492425, Lng: 30.476695298754866},
		{Lat: 59.996694566269, Lng: 30.477467774951155},
		{Lat: 59.985522846219666, Lng: 30.491372346484358},
		{Lat: 59.9734776331996, Lng: 30.54252743681639},
		{Lat: 59.96656200178617, Lng: 30.552827119433577},
		{Lat: 59.9591128611504, Lng: 30.553621053301985},
		{Lat: 59.945472629965536, Lng: 30.540682077014143},
		{Lat: 59.93193933304819, Lng: 30.538150071704084},
		{Lat: 59.92069063807112, Lng: 30.526219606005842},
		{Lat: 59.8887759685014, Lng: 30.5252754684326},
		{Lat: 59.87337726427855, Lng: 30.532571076953108},
		{Lat: 59.86621991030129, Lng: 30.52879452666014},
		{Lat: 59.85465933529358, Lng: 30.503388642871077},
		{Lat: 59.852751939209504, Lng: 30.478669404589827},
		{Lat: 59.847395558876755, Lng: 30.459443330371077},
		{Lat: 59.82596141511424, Lng: 30.4333508010742},
		{Lat: 59.81002495078666, Lng: 30.330747589075262},
		{Lat: 59.82397796859691, Lng: 30.293178893232042},
		{Lat: 59.83578299691027, Lng: 30.28023103212684},
		{Lat: 59.850996682094625, Lng: 30.29092323926955},
		{Lat: 59.87652512736937, Lng: 30.295926020087},
		{Lat: 59.88118594412356, Lng: 30.28675443686291},
		{Lat: 59.88696123875598, Lng: 30.254702824938366},
		{Lat: 59.89260411819843, Lng: 30.247946733331563},
		{Lat: 59.89460722061685, Lng: 30.23787389382699},
		{Lat: 59.90081745317471, Lng: 30.219533717360836},
		{Lat: 59.903951583541954, Lng: 30.21061422347623},
		{Lat: 59.906224690411726, Lng: 30.206501248146314},
		{Lat: 59.90887701706451, Lng: 30.20600005944937},
		{Lat: 59.92205845400411, Lng: 30.211336314284498},
		{Lat: 59.93521315733946, Lng: 30.210664420926268},
		{Lat: 59.946772261294086, Lng: 30.202353596293623},
		{Lat: 59.966258604457494, Lng: 30.216488837802107},
		{Lat: 59.976802045242714, Lng: 30.213457941615278},
		{Lat: 59.98184640764717, Lng: 30.228279828631575},
		{Lat: 60.00888149662998, Lng: 30.23538231810301},
		{Lat: 60.02173941656657, Lng: 30.21937489470213},
		{Lat: 60.03509606030931, Lng: 30.18049359282225},
		{Lat: 60.04054007688507, Lng: 30.157662629687483},
		{Lat: 60.049530432817626, Lng: 30.14880413604022},
	}
	return NewPolygon(points, true)
}

func (p Polygon) MaxLat() float64 {
	var maxLat float64 = -90
	for _, v := range p.polygon.Points() {
		if v.Lat() > maxLat {
			maxLat = v.Lat()
		}
	}
	return maxLat
}

func (p Polygon) MinLat() float64 {
	var minLat float64 = 90
	for _, v := range p.polygon.Points() {
		if v.Lat() < minLat {
			minLat = v.Lat()
		}
	}
	return minLat
}

func (p Polygon) MaxLng() float64 {
	var maxLng float64 = -90
	for _, v := range p.polygon.Points() {
		if v.Lng() > maxLng {
			maxLng = v.Lng()
		}
	}
	return maxLng
}

func (p Polygon) MinLng() float64 {
	var minLng float64 = 90
	for _, v := range p.polygon.Points() {
		if v.Lng() < minLng {
			minLng = v.Lng()
		}
	}
	return minLng
}
