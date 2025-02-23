package main

import (
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Структура для збереження вхідних даних та результату
type Task1Data struct {
	HP, CP, SP, NP, OP, WP, AP float64
	Result                     string
}

// Округлення до 2 знаків після коми
func formatValue(value float64) float64 {
	return math.Round(value*100) / 100
}

// Розрахунок коефіцієнта переходу
func calculateCoefficient(moistureOrCombined float64) float64 {
	return formatValue(100 / (100 - moistureOrCombined))
}

// Розрахунок маси компонента з урахуванням коефіцієнта
func calculateMass(component, coefficient float64) float64 {
	return formatValue(component * coefficient)
}

// Розрахунок нижчої теплоти згоряння для робочої маси (МДж/кг)
func calculateHeat(cp, hp, sp, op, wp float64) float64 {
	qWorking := 339*cp + 1030*hp - 108.8*(op-sp) - 25*wp
	return formatValue(qWorking / 1000)
}

// Розрахунок теплоти згоряння з урахуванням вологи і (для горючої) золи
func calculateHeatMass(baseHeat, moisture float64, ash ...float64) float64 {
	if len(ash) > 0 {
		return formatValue((baseHeat+0.025*moisture)*100 / (100 - moisture - ash[0]))
	}
	return formatValue((baseHeat+0.025*moisture)*100 / (100 - moisture))
}

func task1Handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/task1.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := Task1Data{}

	if r.Method == http.MethodPost {
		// Зчитування даних з форми (якщо поле порожнє, значення буде 0)
		data.HP, _ = strconv.ParseFloat(r.FormValue("hp"), 64)
		data.CP, _ = strconv.ParseFloat(r.FormValue("cp"), 64)
		data.SP, _ = strconv.ParseFloat(r.FormValue("sp"), 64)
		data.NP, _ = strconv.ParseFloat(r.FormValue("np"), 64)
		data.OP, _ = strconv.ParseFloat(r.FormValue("op"), 64)
		data.WP, _ = strconv.ParseFloat(r.FormValue("wp"), 64)
		data.AP, _ = strconv.ParseFloat(r.FormValue("ap"), 64)

		// Розрахунок коефіцієнтів
		coefficientDry := calculateCoefficient(data.WP)
		coefficientCombustible := calculateCoefficient(data.WP + data.AP)

		// Розрахунок компонентів для сухої маси
		hpDry := calculateMass(data.HP, coefficientDry)
		cpDry := calculateMass(data.CP, coefficientDry)
		spDry := calculateMass(data.SP, coefficientDry)
		npDry := calculateMass(data.NP, coefficientDry)
		opDry := calculateMass(data.OP, coefficientDry)
		apDry := calculateMass(data.AP, coefficientDry)

		// Розрахунок компонентів для горючої маси
		hpComb := calculateMass(data.HP, coefficientCombustible)
		cpComb := calculateMass(data.CP, coefficientCombustible)
		spComb := calculateMass(data.SP, coefficientCombustible)
		npComb := calculateMass(data.NP, coefficientCombustible)
		opComb := calculateMass(data.OP, coefficientCombustible)

		// Розрахунок нижчої теплоти згоряння для робочої маси
		heatWorking := calculateHeat(data.CP, data.HP, data.SP, data.OP, data.WP)
		// Теплота для сухої маси
		heatDry := calculateHeatMass(heatWorking, data.WP)
		// Теплота для горючої маси (враховуючи золу)
		heatCombustible := calculateHeatMass(heatWorking, data.WP, data.AP)

		// Формування результатного рядка
		data.Result = "Вхідні дані:\n" +
			"HP: " + strconv.FormatFloat(data.HP, 'f', 2, 64) + "%, " +
			"CP: " + strconv.FormatFloat(data.CP, 'f', 2, 64) + "%, " +
			"SP: " + strconv.FormatFloat(data.SP, 'f', 2, 64) + "%, " +
			"NP: " + strconv.FormatFloat(data.NP, 'f', 2, 64) + "%, " +
			"OP: " + strconv.FormatFloat(data.OP, 'f', 2, 64) + "%, " +
			"WP: " + strconv.FormatFloat(data.WP, 'f', 2, 64) + "%, " +
			"AP: " + strconv.FormatFloat(data.AP, 'f', 2, 64) + "%\n\n" +

			"Коефіцієнти переходу:\n" +
			"Робоча -> суха: " + strconv.FormatFloat(coefficientDry, 'f', 2, 64) + "\n" +
			"Робоча -> горюча: " + strconv.FormatFloat(coefficientCombustible, 'f', 2, 64) + "\n\n" +

			"Склад сухої маси:\n" +
			"HP: " + strconv.FormatFloat(hpDry, 'f', 2, 64) + "%, " +
			"CP: " + strconv.FormatFloat(cpDry, 'f', 2, 64) + "%, " +
			"SP: " + strconv.FormatFloat(spDry, 'f', 2, 64) + "%, " +
			"NP: " + strconv.FormatFloat(npDry, 'f', 2, 64) + "%, " +
			"OP: " + strconv.FormatFloat(opDry, 'f', 2, 64) + "%, " +
			"AP: " + strconv.FormatFloat(apDry, 'f', 2, 64) + "%\n\n" +

			"Склад горючої маси:\n" +
			"HP: " + strconv.FormatFloat(hpComb, 'f', 2, 64) + "%, " +
			"CP: " + strconv.FormatFloat(cpComb, 'f', 2, 64) + "%, " +
			"SP: " + strconv.FormatFloat(spComb, 'f', 2, 64) + "%, " +
			"NP: " + strconv.FormatFloat(npComb, 'f', 2, 64) + "%, " +
			"OP: " + strconv.FormatFloat(opComb, 'f', 2, 64) + "%\n\n" +

			"Нижча теплота згоряння:\n" +
			"Робоча маса: " + strconv.FormatFloat(heatWorking, 'f', 2, 64) + " МДж/кг\n" +
			"Суха маса: " + strconv.FormatFloat(heatDry, 'f', 2, 64) + " МДж/кг\n" +
			"Горюча маса: " + strconv.FormatFloat(heatCombustible, 'f', 2, 64) + " МДж/кг"
	}

	tmpl.Execute(w, data)
}


// Структура для даних другого калькулятора
type Task2Data struct {
	Carbon, Hydrogen, Oxygen, Sulfur, OilHeat, FuelMoisture, Ash, Vanadium float64
	Result                                                                  string
}

// Обробник для другого калькулятора
func task2Handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/task2.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := Task2Data{}

	if r.Method == http.MethodPost {
		data.Carbon, _ = strconv.ParseFloat(r.FormValue("carbon"), 64)
		data.Hydrogen, _ = strconv.ParseFloat(r.FormValue("hydrogen"), 64)
		data.Oxygen, _ = strconv.ParseFloat(r.FormValue("oxygen"), 64)
		data.Sulfur, _ = strconv.ParseFloat(r.FormValue("sulfur"), 64)
		data.OilHeat, _ = strconv.ParseFloat(r.FormValue("oilHeat"), 64)
		data.FuelMoisture, _ = strconv.ParseFloat(r.FormValue("fuelMoisture"), 64)
		data.Ash, _ = strconv.ParseFloat(r.FormValue("ash"), 64)
		data.Vanadium, _ = strconv.ParseFloat(r.FormValue("vanadium"), 64)

		// Обчислення множників згідно із завданням
		factor1 := (100 - data.FuelMoisture - data.Ash) / 100
		factor2 := (100 - data.FuelMoisture/10 - data.Ash/10) / 100
		factor3 := (100 - data.FuelMoisture) / 100

		// Перерахунок компонентів для робочої маси
		carbonWorking := formatValue(data.Carbon * factor1)
		hydrogenWorking := formatValue(data.Hydrogen * factor1)
		oxygenWorking := formatValue(data.Oxygen * factor2)
		sulfurWorking := formatValue(data.Sulfur * factor1)
		ashWorking := formatValue(data.Ash * factor3)
		vanadiumWorking := formatValue(data.Vanadium * factor3)

		// Перерахунок нижчої теплоти згоряння для робочої маси
		lowerHeat := formatValue(data.OilHeat*factor1 - 0.025*data.FuelMoisture)

		// Формування рядка з результатами
		data.Result = "Вхідні дані:\n" +
			"Вуглець: " + strconv.FormatFloat(data.Carbon, 'f', 2, 64) + "%, " +
			"Водень: " + strconv.FormatFloat(data.Hydrogen, 'f', 2, 64) + "%, " +
			"Кисень: " + strconv.FormatFloat(data.Oxygen, 'f', 2, 64) + "%, " +
			"Сірка: " + strconv.FormatFloat(data.Sulfur, 'f', 2, 64) + "%,\n" +
			"Нижча теплота горючої маси: " + strconv.FormatFloat(data.OilHeat, 'f', 2, 64) + " МДж/кг, " +
			"Вологість: " + strconv.FormatFloat(data.FuelMoisture, 'f', 2, 64) + "%, " +
			"Зольність: " + strconv.FormatFloat(data.Ash, 'f', 2, 64) + "%,\n" +
			"Вміст ванадію: " + strconv.FormatFloat(data.Vanadium, 'f', 2, 64) + " мг/кг\n\n" +
			"Склад робочої маси мазуту:\n" +
			"C: " + strconv.FormatFloat(carbonWorking, 'f', 2, 64) + "%, " +
			"H: " + strconv.FormatFloat(hydrogenWorking, 'f', 2, 64) + "%, " +
			"O: " + strconv.FormatFloat(oxygenWorking, 'f', 2, 64) + "%,\n" +
			"S: " + strconv.FormatFloat(sulfurWorking, 'f', 2, 64) + "%, " +
			"A: " + strconv.FormatFloat(ashWorking, 'f', 2, 64) + "%, " +
			"V: " + strconv.FormatFloat(vanadiumWorking, 'f', 2, 64) + " мг/кг\n\n" +
			"Нижча теплота згоряння (робоча маса): " + strconv.FormatFloat(lowerHeat, 'f', 2, 64) + " МДж/кг"
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
