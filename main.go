package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Cell struct {
	ClassList []string `json:"classList"`
	X         string   `json:"x"`
	Y         string   `json:"y"`
}

var lastMove []Cell
var filePath = "last_move.txt"

func sendMove(w http.ResponseWriter, r *http.Request) {
	var move []Cell
	err := json.NewDecoder(r.Body).Decode(&move)
	if err != nil {
		log.Println("Erro ao processar a solicitação:", err)
		http.Error(w, "Invalid move data", http.StatusBadRequest)
		return
	}

	lastMove = move
	saveToFile(move)

	w.WriteHeader(http.StatusOK)
}

func getMove(w http.ResponseWriter, r *http.Request) {
	loadFromFile(&lastMove)

	if len(lastMove) == 0 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "no move"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "move available", "move": lastMove})

	lastMove = nil
}

func saveToFile(move []Cell) {
	moveBytes, err := json.Marshal(move)
	if err != nil {
		log.Println("Erro ao serializar o movimento:", err)
		return
	}

	err = ioutil.WriteFile(filePath, moveBytes, 0644)
	if err != nil {
		log.Println("Erro ao escrever o arquivo:", err)
		return
	}
}

func loadFromFile(move *[]Cell) {
	moveBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("Erro ao ler o arquivo:", err)
		return
	}

	if len(moveBytes) == 0 {
		log.Println("Arquivo vazio. Nenhum movimento disponível.")
		return
	}

	err = json.Unmarshal(moveBytes, move)
	if err != nil {
		log.Println("Erro ao desserializar o movimento:", err)
		return
	}
}

func main() {
	r := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:5500"},
		AllowedHeaders: []string{"Content-Type"},
	})

	r.Use(c.Handler)

	r.HandleFunc("/api/send-move", sendMove).Methods("POST")
	r.HandleFunc("/api/check", getMove).Methods("GET")

	log.Fatal(http.ListenAndServe(":5003", r))
}
