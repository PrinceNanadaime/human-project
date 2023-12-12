package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ApiServer struct {
	human       HumanFactService
	calculation CalculateService
	database    DataBaseService
}

func NewApiServer(human HumanFactService, calculation CalculateService, database DataBaseService) ApiServer {
	return ApiServer{
		human:       human,
		calculation: calculation,
		database:    database,
	}
}

func (s *ApiServer) Start() error {
	http.HandleFunc("/human-fact", s.HandleGetHumanFact)
	http.HandleFunc("/calculate-length", s.HandleCalculation)
	http.HandleFunc("/db/insert", s.HandleInsertionToDatabase)
	http.HandleFunc("/db/get-all", s.HandleGetAllFromDatabase)
	return http.ListenAndServe("localhost:8080", nil)
}

func (s *ApiServer) HandleGetHumanFact(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на получение нового факта о человеке!")
	fact, err := s.human.CreateHumanFact()
	if err != nil {
		log.Println("В процессе получения нового факта поймали исключение:", err.Error())
		WriteJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJson(w, http.StatusOK, fact)
	log.Println("Запрос на получение нового факта о человеке успешно выполнен!")
}

func (s *ApiServer) HandleCalculation(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на вычисление длины факта о человеке!")
	body := HumanFact{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println("В процессе вычисления длины факта поймали исключение:", err.Error())
		WriteJson(w, http.StatusInternalServerError, err.Error())
	}
	fact := s.calculation.CalculateLength(body)
	WriteJson(w, http.StatusOK, fact)
	log.Println("Запрос на получение вычисление длины факта о человеке успешно выполнен!")
}

func (s *ApiServer) HandleInsertionToDatabase(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на запись нового факта о человеке в БД!")
	body := HumanFact{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println("В процессе записи данных в БД поймали исключение:", err.Error())
		WriteJson(w, http.StatusInternalServerError, err.Error())
	}

	fact, err := s.database.Insert(body)
	if err != nil {
		log.Println("В процессе записи данных в БД поймали исключение:", err.Error())
		WriteJson(w, http.StatusInternalServerError, err.Error())
	}
	WriteJson(w, http.StatusOK, fact)
	log.Println("Запрос на запись нового факта о человеке в БД успешно выполнен!")
}

func (s *ApiServer) HandleGetAllFromDatabase(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на вывод всех фактов о человеке из БД!")
	fact, err := s.database.GetAll()
	if err != nil {
		log.Println("В процессе вывода всех данных в БД поймали исключение:", err.Error())
		WriteJson(w, http.StatusInternalServerError, err.Error())
	}
	WriteJson(w, http.StatusOK, fact)
	log.Println("Запрос на вывод всех фактов о человеке из БД успешно выполнен!")
}

func WriteJson(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	return encoder.Encode(v)
}
