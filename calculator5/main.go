package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

type TaskData struct {
	Result string
}

func calculateReliabilitySingleLineSystem() float64 {
	failureRates := []float64{0.01, 0.07, 0.015, 0.02, 0.03 * 6}
	failureRateSum := 0.0
	for _, rate := range failureRates {
		failureRateSum += rate
	}
	meanRecoveryTime := (0.01*30 + 0.07*10 + 0.015*100 + 0.02*15 + 0.18*2) / failureRateSum
	return (failureRateSum * meanRecoveryTime) / 8760
}

func calculateReliabilityDoubleLineSystem() float64 {
	failureRateSingleLine := 0.295
	failureRateTwoLinesSimultaneous := 2*failureRateSingleLine*(13.6*1e-4) + 5.89*1e-7
	return failureRateTwoLinesSimultaneous + 0.02
}

func calculatePowerLoss(emergencyRate float64, plannedRate float64) float64 {
	const emergencyPowerLoss = 14900.0
	const plannedPowerLoss = 132400.0

	return (emergencyRate * emergencyPowerLoss) + (plannedRate * plannedPowerLoss)
}

func task1Handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("task1.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := TaskData{}

	if r.Method == http.MethodPost {
		singleLineReliability := calculateReliabilitySingleLineSystem()
		doubleLineReliability := calculateReliabilityDoubleLineSystem()

		moreReliableSystem := ""
		if singleLineReliability > doubleLineReliability {
			moreReliableSystem = "Одноколова система більш надійна."
		} else {
			moreReliableSystem = "Двоколова система більш надійна."
		}

		data.Result = fmt.Sprintf(
			"Надійність одноколової системи: %.6f\nНадійність двоколової системи: %.6f\n%s",
			singleLineReliability, doubleLineReliability, moreReliableSystem,
		)
	}
	tmpl.Execute(w, data)
}

func task2Handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("task2.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := TaskData{}

	if r.Method == http.MethodPost {
		emergencyRateStr := r.FormValue("emergencyRate")
		plannedRateStr := r.FormValue("plannedRate")

		emergencyRate, err1 := strconv.ParseFloat(emergencyRateStr, 64)
		plannedRate, err2 := strconv.ParseFloat(plannedRateStr, 64)

		if err1 == nil && err2 == nil {
			result := calculatePowerLoss(emergencyRate, plannedRate)
			data.Result = fmt.Sprintf("Збитки: %.2f грн", result)
		} else {
			data.Result = "Помилка: введіть коректні значення питомих збитків."
		}
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/task1", task1Handler)
	http.HandleFunc("/task2", task2Handler)

	log.Println("Сервер запущено на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
