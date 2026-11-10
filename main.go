package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Appointment struct {
	ID             int     `json:"appointment_id"`
	ClientName     string  `json:"client_name"`
	SpecialistName string  `json:"specialist_name"`
	ServiceName    string  `json:"service_name"`
	Date           string  `json:"date"`
	Time           string  `json:"time"`
	Status         string  `json:"status"`
	Price          float64 `json:"price"`
}

var db *sql.DB

func initDB() {
	var err error
	dsn := "root:89236488187@tcp(127.0.0.1:3306)/spa_booking?parseTime=true"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic("Не удалось подключиться к БД. Проверь пароль. Ошибка: " + err.Error())
	}
	fmt.Println("Успешное подключение к базе данных MySQL!")
}

func enableCors(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}

func getAppointments(w http.ResponseWriter, r *http.Request) {
	if enableCors(w, r) {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	query := `
		SELECT 
			a.appointment_id, 
			c.first_name AS client_name, 
			s.last_name AS specialist_name, 
			srv.name AS service_name, 
			a.appointment_date, 
			a.appointment_time, 
			a.status,
			srv.price
		FROM Appointment a
		JOIN Client c ON a.client_id = c.client_id
		JOIN Specialist s ON a.specialist_id = s.specialist_id
		JOIN Service srv ON a.service_id = srv.service_id
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var result []Appointment
	for rows.Next() {
		var a Appointment
		err = rows.Scan(&a.ID, &a.ClientName, &a.SpecialistName, &a.ServiceName, &a.Date, &a.Time, &a.Status, &a.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result = append(result, a)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/api/bookings", getAppointments)

	port := ":8080"
	fmt.Printf("Сервер SPA BOOKING запущен на http://localhost%s\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
