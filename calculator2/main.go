package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type EmissionResult struct {
	CoalEmissionIndex    float64
	CoalEmission         float64
	FuelOilEmissionIndex float64
	FuelOilEmission      float64
}

var tmpl, err = template.ParseFiles("index.html")

func calculateEmissions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	coalCombustion, _ := strconv.ParseFloat(r.FormValue("coalCombustion"), 64)
	fuelOilCombustion, _ := strconv.ParseFloat(r.FormValue("fuelOilCombustion"), 64)
	coalAshContent, _ := strconv.ParseFloat(r.FormValue("coalAshContent"), 64)
	fuelOilAshContent, _ := strconv.ParseFloat(r.FormValue("fuelOilAshContent"), 64)
	coalFuelMass, _ := strconv.ParseFloat(r.FormValue("coalFuelMass"), 64)
	fuelOilFuelMass, _ := strconv.ParseFloat(r.FormValue("fuelOilFuelMass"), 64)
	dustRemovalEfficiency, _ := strconv.ParseFloat(r.FormValue("dustRemovalEfficiency"), 64)

	coalEmissionIndex, coalEmission := calculateCoalEmissions(coalCombustion, coalAshContent, coalFuelMass, dustRemovalEfficiency)
	fuelOilEmissionIndex, fuelOilEmission := calculateFuelOilEmissions(fuelOilCombustion, fuelOilAshContent, fuelOilFuelMass, dustRemovalEfficiency)

	result := EmissionResult{
		CoalEmissionIndex:    coalEmissionIndex,
		CoalEmission:         coalEmission,
		FuelOilEmissionIndex: fuelOilEmissionIndex,
		FuelOilEmission:      fuelOilEmission,
	}
	tmpl.Execute(w, result)
}

func calculateCoalEmissions(coalCombustion, coalAshContent, coalFuelMass, dustRemovalEfficiency float64) (float64, float64) {
	emissionIndex := math.Pow(10, 6) / coalCombustion * 0.8 * coalAshContent / (100 - 1.5) * (1 - dustRemovalEfficiency)
	emission := math.Pow(10, -6) * emissionIndex * coalCombustion * coalFuelMass
	return math.Round(emissionIndex*100) / 100, math.Round(emission*100) / 100
}

func calculateFuelOilEmissions(fuelOilCombustion, fuelOilAshContent, fuelOilFuelMass, dustRemovalEfficiency float64) (float64, float64) {
	emissionIndex := math.Pow(10, 6) / fuelOilCombustion * 1 * fuelOilAshContent / (100 - 1.5) * (1 - dustRemovalEfficiency)
	emission := math.Pow(10, -6) * emissionIndex * fuelOilCombustion * fuelOilFuelMass
	return math.Round(emissionIndex*100) / 100, math.Round(emission*100) / 100
}

func main() {
	http.HandleFunc("/", calculateEmissions)
	fmt.Println("Сервер запущено на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
