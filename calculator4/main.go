package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
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

func calculateCableParameters(voltage float64) (float64, float64) {
	const shortCircuitCurrent = 2500.0
	const disconnectionTime = 2.5
	const loadPower = 1300.0
	const economicDensity = 1.4
	const thermalConstant = 92.0

	inNormal := loadPower / (2 * math.Sqrt(3) * voltage)
	economicSection := inNormal / economicDensity
	thermalSection := shortCircuitCurrent * math.Sqrt(disconnectionTime) / thermalConstant

	return economicSection, thermalSection
}


func calculateShortCircuitCurrent(pKz float64) float64 {
	const voltage = 10000.0
	const resistanceK1 = 0.55
	const resistanceK2 = 1.84
	const sqrt3 = 1.732

	return pKz / (voltage * sqrt3 * (resistanceK1 + resistanceK2))
}

func calculateShortCircuitCurrents(Rh, Xh, Rm, Xm float64) string {
	const nominalVoltage = 115.0
	const baseVoltage = 11.0
	const sqrt3 = 1.732
	const multiplier = 1000.0

	Xt := (11.1 * math.Pow(nominalVoltage, 2)) / (100 * 6.3)

	Zsh := math.Sqrt(math.Pow(Rh, 2) + math.Pow(Xh+Xt, 2))
	ZshMin := math.Sqrt(math.Pow(Rm, 2) + math.Pow(Xm+Xt, 2))

	Ish3Normal := (nominalVoltage * multiplier) / (sqrt3 * Zsh)
	Ish3Min := (nominalVoltage * multiplier) / (sqrt3 * ZshMin)

	Ish2Normal := Ish3Normal * (sqrt3 / 2)
	Ish2Min := Ish3Min * (sqrt3 / 2)

	k := math.Pow(baseVoltage, 2) / math.Pow(nominalVoltage, 2)
	ZshTrue := math.Sqrt(math.Pow(Rh*k, 2) + math.Pow((Xh+Xt)*k, 2))
	ZshMinTrue := math.Sqrt(math.Pow(Rm*k, 2) + math.Pow((Xm+Xt)*k, 2))

	DIsh3Normal := (baseVoltage * multiplier) / (sqrt3 * ZshTrue)
	DIsh3Min := (baseVoltage * multiplier) / (sqrt3 * ZshMinTrue)

	DIsh2Normal := DIsh3Normal * (sqrt3 / 2)
	DIsh2Min := DIsh3Min * (sqrt3 / 2)

	return fmt.Sprintf(`
Результати розрахунків:
Струм трифазного КЗ (приведений до 110 кВ):
Нормальний режим: %.2f А
Мінімальний режим: %.2f А

Струм двофазного КЗ (приведений до 110 кВ):
Нормальний режим: %.2f А
Мінімальний режим: %.2f А

Дійсний струм трифазного КЗ (на шинах 10 кВ):
Нормальний режим: %.2f А
Мінімальний режим: %.2f А

Дійсний струм двофазного КЗ (на шинах 10 кВ):
Нормальний режим: %.2f А
Мінімальний режим: %.2f А

Аварійний режим на даній підстанції не передбачений.
`, Ish3Normal, Ish3Min, Ish2Normal, Ish2Min, DIsh3Normal, DIsh3Min, DIsh2Normal, DIsh2Min)
}


func task1Handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("task1.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := TaskData{}

	if r.Method == http.MethodPost {
		voltageStr := r.FormValue("voltage")
		voltage, err := strconv.ParseFloat(voltageStr, 64)
		if err == nil {
			economicSection, thermalSection := calculateCableParameters(voltage)
			data.Result = "Економічний переріз: " + strconv.FormatFloat(economicSection, 'f', 2, 64) + " мм²\n" +
				"Термічний переріз: " + strconv.FormatFloat(thermalSection, 'f', 2, 64) + " мм²"
		} else {
			data.Result = "Помилка: введіть коректне значення напруги."
		}
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
		pKzStr := r.FormValue("pKz")
		pKz, err := strconv.ParseFloat(pKzStr, 64)
		if err == nil {
			current := calculateShortCircuitCurrent(pKz)
			data.Result = "Струм короткого замикання: " + strconv.FormatFloat(current, 'f', 8, 64) + " кА"
		} else {
			data.Result = "Помилка: введіть коректне значення потужності КЗ."
		}
	}

	tmpl.Execute(w, data)
}

func task3Handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("task3.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := TaskData{}

	if r.Method == http.MethodPost {
		resistanceHighStr := r.FormValue("resistanceHigh")
		reactanceHighStr := r.FormValue("reactanceHigh")
		resistanceMediumStr := r.FormValue("resistanceMedium")
		reactanceMediumStr := r.FormValue("reactanceMedium")

		resistanceHigh, err1 := strconv.ParseFloat(resistanceHighStr, 64)
		reactanceHigh, err2 := strconv.ParseFloat(reactanceHighStr, 64)
		resistanceMedium, err3 := strconv.ParseFloat(resistanceMediumStr, 64)
		reactanceMedium, err4 := strconv.ParseFloat(reactanceMediumStr, 64)

		if err1 == nil && err2 == nil && err3 == nil && err4 == nil {
			data.Result = calculateShortCircuitCurrents(resistanceHigh, reactanceHigh, resistanceMedium, reactanceMedium)
		} else {
			data.Result = "Будь ласка, введіть всі значення коректно."
		}
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/task1", task1Handler)
	http.HandleFunc("/task2", task2Handler)
	http.HandleFunc("/task3", task3Handler)

	log.Println("Сервер запущено на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
