package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

// Equipment represents a piece of equipment in the calculator
type Equipment struct {
	Name                string  `json:"name"`
	Efficiency          float64 `json:"efficiency"`
	PowerFactor         float64 `json:"powerFactor"`
	Voltage             float64 `json:"voltage"`
	Quantity            float64 `json:"quantity"`
	NominalPower        float64 `json:"nominalPower"`
	UtilizationFactor   float64 `json:"utilizationFactor"`
	ReactivePowerFactor float64 `json:"reactivePowerFactor"`
	TotalPower          float64 `json:"totalPower"`
	Current             float64 `json:"current"`
}

// PowerCalculations represents the power calculation results
type PowerCalculations struct {
	ActiveLoad   float64 `json:"activeLoad"`
	ReactiveLoad float64 `json:"reactiveLoad"`
	FullPower    float64 `json:"fullPower"`
	GroupCurrent float64 `json:"groupCurrent"`
}

// DepartmentCalculations represents department-specific calculations
type DepartmentCalculations struct {
	UtilizationCoefficient float64 `json:"utilizationCoefficient"`
	EffectiveAmount        float64 `json:"effectiveAmount"`
}

// TransformerCalculations represents transformer-specific calculations
type TransformerCalculations struct {
	ActiveLoad   float64 `json:"activeLoad"`
	ReactiveLoad float64 `json:"reactiveLoad"`
	FullPower    float64 `json:"fullPower"`
	GroupCurrent float64 `json:"groupCurrent"`
}

// CalculationState represents the overall calculation state
type CalculationState struct {
	GroupUtilizationCoefficient float64                 `json:"groupUtilizationCoefficient"`
	EffectiveEpAmount           float64                 `json:"effectiveEpAmount"`
	CalculatedPower             PowerCalculations       `json:"calculatedPower"`
	DepartmentCalculations      DepartmentCalculations  `json:"departmentCalculations"`
	TransformerCalculations     TransformerCalculations `json:"transformerCalculations"`
}

// CalculationRequest represents the input for calculations
type CalculationRequest struct {
	EquipmentList          []Equipment `json:"equipmentList"`
	ActiveCoefficient      float64     `json:"activeCoefficient"`
	TransformerCoefficient float64     `json:"transformerCoefficient"`
}

// GetInitialEquipmentList returns the initial equipment list
func GetInitialEquipmentList() []Equipment {
	return []Equipment{
		{
			Name:                "Шліфувальний верстат",
			Efficiency:          0.92,
			PowerFactor:         0.9,
			Voltage:             0.38,
			Quantity:            4,
			NominalPower:        20,
			UtilizationFactor:   0.15,
			ReactivePowerFactor: 1.33,
		},
		{
			Name:                "Свердлильний верстат",
			Efficiency:          0.92,
			PowerFactor:         0.9,
			Voltage:             0.38,
			Quantity:            2,
			NominalPower:        14,
			UtilizationFactor:   0.12,
			ReactivePowerFactor: 1,
		},
		{
			Name:                "Фугувальний верстат",
			Efficiency:          0.92,
			PowerFactor:         0.9,
			Voltage:             0.38,
			Quantity:            4,
			NominalPower:        42,
			UtilizationFactor:   0.15,
			ReactivePowerFactor: 1.33,
		},
		{
			Name:                "Циркулярна пила",
			Efficiency:          0.92,
			PowerFactor:         0.9,
			Voltage:             0.38,
			Quantity:            1,
			NominalPower:        36,
			UtilizationFactor:   0.3,
			ReactivePowerFactor: 1.55,
		},
		{
			Name:                "Прес",
			Efficiency:          0.92,
			PowerFactor:         0.9,
			Voltage:             0.38,
			Quantity:            1,
			NominalPower:        20,
			UtilizationFactor:   0.5,
			ReactivePowerFactor: 0.75,
		},
		{
			Name:                "Полірувальний верстат",
			Efficiency:          0.92,
			PowerFactor:         0.9,
			Voltage:             0.38,
			Quantity:            1,
			NominalPower:        40,
			UtilizationFactor:   0.21,
			ReactivePowerFactor: 1,
		},
		{
			Name:                "Фрезерний верстат",
			Efficiency:          0.92,
			PowerFactor:         0.9,
			Voltage:             0.38,
			Quantity:            2,
			NominalPower:        32,
			UtilizationFactor:   0.2,
			ReactivePowerFactor: 1,
		},
		{
			Name:                "Вентилятор",
			Efficiency:          0.92,
			PowerFactor:         0.9,
			Voltage:             0.38,
			Quantity:            1,
			NominalPower:        20,
			UtilizationFactor:   0.65,
			ReactivePowerFactor: 0.75,
		},
	}
}

// CalculateEquipmentParameters calculates parameters for each equipment
func CalculateEquipmentParameters(equipmentList []Equipment, activeCoefficient, transformerCoefficient float64) CalculationState {
	var sumOfPowerProduct float64
	var sumOfPowerUtilProduct float64
	var sumOfSquaredPowerProduct float64

	for i := range equipmentList {
		quantity := equipmentList[i].Quantity
		power := equipmentList[i].NominalPower
		totalPower := quantity * power

		equipmentList[i].TotalPower = totalPower
		equipmentList[i].Current = calculateCurrent(
			totalPower,
			equipmentList[i].Voltage,
			equipmentList[i].PowerFactor,
			equipmentList[i].Efficiency,
		)

		sumOfPowerProduct += totalPower
		sumOfPowerUtilProduct += totalPower * equipmentList[i].UtilizationFactor
		sumOfSquaredPowerProduct += quantity * power * power
	}

	return updateCalculationState(sumOfPowerProduct, sumOfPowerUtilProduct, sumOfSquaredPowerProduct, activeCoefficient, transformerCoefficient)
}

// calculateCurrent calculates the current for an equipment
func calculateCurrent(power, voltage, powerFactor, efficiency float64) float64 {
	return roundToTwoDecimalPlaces(power / (math.Sqrt(3.0) * voltage * powerFactor * efficiency))
}

// updateCalculationState updates the calculation state
func updateCalculationState(
	totalPower,
	totalPowerUtil,
	totalSquaredPower,
	activeCoef,
	transformerCoef float64,
) CalculationState {
	groupUtilCoef := roundToTwoDecimalPlaces(totalPowerUtil / totalPower)
	effectiveAmount := roundToTwoDecimalPlaces((totalPower * totalPower) / totalSquaredPower)

	const voltage = 0.38
	const nominalPower = 23.0
	const tangentPhi = 1.58

	activeLoad := roundToTwoDecimalPlaces(activeCoef * totalPowerUtil)
	reactiveLoad := roundToTwoDecimalPlaces(groupUtilCoef * nominalPower * tangentPhi)
	fullPower := roundToTwoDecimalPlaces(math.Sqrt(activeLoad*activeLoad + reactiveLoad*reactiveLoad))
	groupCurrent := roundToTwoDecimalPlaces(activeLoad / voltage)

	calculatedPower := PowerCalculations{
		ActiveLoad:   activeLoad,
		ReactiveLoad: reactiveLoad,
		FullPower:    fullPower,
		GroupCurrent: groupCurrent,
	}

	departmentCalculations := DepartmentCalculations{
		UtilizationCoefficient: 752.0 / 2330.0,
		EffectiveAmount:        2330.0 * 2330.0 / 96399.0,
	}

	transformerCalculations := calculateTransformerParameters(transformerCoef)

	return CalculationState{
		GroupUtilizationCoefficient: groupUtilCoef,
		EffectiveEpAmount:           effectiveAmount,
		CalculatedPower:             calculatedPower,
		DepartmentCalculations:      departmentCalculations,
		TransformerCalculations:     transformerCalculations,
	}
}

// calculateTransformerParameters calculates transformer-specific parameters
func calculateTransformerParameters(coefficient float64) TransformerCalculations {
	activeLoad := roundToTwoDecimalPlaces(coefficient * 752.0)
	reactiveLoad := roundToTwoDecimalPlaces(coefficient * 657.0)
	fullPower := roundToTwoDecimalPlaces(math.Sqrt(activeLoad*activeLoad + reactiveLoad*reactiveLoad))
	groupCurrent := roundToTwoDecimalPlaces(activeLoad / 0.38)

	return TransformerCalculations{
		ActiveLoad:   activeLoad,
		ReactiveLoad: reactiveLoad,
		FullPower:    fullPower,
		GroupCurrent: groupCurrent,
	}
}

// roundToTwoDecimalPlaces rounds a number to two decimal places
func roundToTwoDecimalPlaces(value float64) float64 {
	return math.Round(value*100) / 100
}

// indexHandler handles the root endpoint
func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// initialDataHandler provides the initial equipment list
func initialDataHandler(w http.ResponseWriter, r *http.Request) {
	equipmentList := GetInitialEquipmentList()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(equipmentList)
}

// calculateHandler performs the calculations based on the input
func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var request CalculationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result := CalculateEquipmentParameters(request.EquipmentList, request.ActiveCoefficient, request.TransformerCoefficient)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	// Handle static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handle routes
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/initial-data", initialDataHandler)
	http.HandleFunc("/api/calculate", calculateHandler)

	// Start server
	port := 8080
	fmt.Printf("Server is running on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
