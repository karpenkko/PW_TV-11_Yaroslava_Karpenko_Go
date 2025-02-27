package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"html/template"
)

// обчислює ймовірність нормального розподілу
func normalPDF(x, mean, stdDev float64) float64 {
	return (1 / (stdDev * math.Sqrt(2*math.Pi))) *
		math.Exp(-0.5*math.Pow((x-mean)/stdDev, 2))
}

// обчислює ймовірність у заданому діапазоні
func calculateNormalProbability(mean, stdDev, lower, upper float64) float64 {
	steps := 1000
	stepSize := (upper - lower) / float64(steps)
	area := 0.0
	for i := 0; i < steps; i++ {
		x1 := lower + float64(i)*stepSize
		x2 := x1 + stepSize
		area += (normalPDF(x1, mean, stdDev) + normalPDF(x2, mean, stdDev)) / 2 * stepSize
	}
	return area
}

// структура для передавання результату в шаблон
type CalculationResult struct {
	ProfitImproved float64
}

// Обробник для головної сторінки
func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Не вдалося завантажити сторінку", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		power, _ := strconv.ParseFloat(r.FormValue("power"), 64)
		deviationImproved, _ := strconv.ParseFloat(r.FormValue("deviationImproved"), 64)
		cost, _ := strconv.ParseFloat(r.FormValue("cost"), 64)

		probabilityNoImbalance2 := calculateNormalProbability(5.0, deviationImproved, 4.75, 5.25)

		w2NoImbalance := power * 24 * probabilityNoImbalance2
		w2Imbalance := power * 24 * (1 - probabilityNoImbalance2)

		profitImproved := math.Round(((w2NoImbalance * cost) - (w2Imbalance * cost)) * 100) / 100

		result := CalculationResult{ProfitImproved: profitImproved}
		tmpl.Execute(w, result)
	} else {
		tmpl.Execute(w, nil)
	}
}

func main() {
	http.HandleFunc("/", handleIndex)
	fmt.Println("Сервер запущено на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

