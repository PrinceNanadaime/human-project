package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
)

type DataBaseService interface {
	Insert(fact HumanFact) (bool, error)
	GetAll() ([]HumanFact, error)
}

type DataBase struct {
	db *sql.DB
}

type HumanFactDTO struct {
	Id     int
	Fact   string `json:"fact"`
	Length int    `json:"length"`
}

func NewDataBaseService(db *sql.DB, err error) DataBaseService {
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS fact (fact_id INTEGER PRIMARY KEY, fact_desc varchar(200) NOT NULL, fact_length INTEGER NOT NULL)")
	if err != nil {
		panic(err)
	}
	log.Println("Подключение к БД выполнено!")
	return &DataBase{
		db: db,
	}
}

func (d *DataBase) Insert(fact HumanFact) (bool, error) {
	log.Println("Начало записи в БД факта: ", fact.Fact)
	_, err := d.db.Exec("INSERT INTO fact values ($1, $2, $3)", rand.Intn(1000000), fact.Fact, fact.Length)
	if err != nil {
		log.Println("В процессе записи факта в БД поймали исключение", err.Error())
		return false, err
	}
	log.Println("Успешный конец записи в БД факта: ", fact.Fact)
	return true, err
}

func (d *DataBase) GetAll() ([]HumanFact, error) {
	log.Println("Начало поиска всех объектов типа HumanFact")
	facts, err := d.db.Query("SELECT * from fact")
	if err != nil {
		log.Println("В процессе поиска поймали исключение", err.Error())
		panic(err)
	}
	var foundFacts []HumanFact
	for facts.Next() {
		var h HumanFactDTO
		err = facts.Scan(&h.Id, &h.Fact, &h.Length)
		if err != nil {
			log.Println("В процессе десериализации объекта поймали исключение", err.Error())
			panic(err)
		}
		foundFacts = append(foundFacts, HumanFact{
			Fact:   h.Fact,
			Length: h.Length})
	}
	log.Println("Успешный конец поиска всех объектов типа HumanFact")
	return foundFacts, err
}
