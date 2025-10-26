package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Client representa a estrutura de dados do cliente
type Client struct {
	ID         string `json:"id"`
	Nome       string `json:"nome"`
	Nascimento string `json:"nascimento"`
	Endereco   string `json:"endereco"`
	Telefone   string `json:"telefone"`
}

var clients = make(map[string]Client)
var currentID = 1 // Começa em 1 pois pré-populamos um cliente

// respondWithError envia uma resposta de erro JSON
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON envia uma resposta JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		err := json.NewEncoder(w).Encode(payload)
		if err != nil {
			log.Printf("Erro ao escrever resposta JSON: %v", err)
		}
	}
}

// createClient cria um novo cliente
func createClient(w http.ResponseWriter, r *http.Request) {
	var client Client
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Requisição com JSON inválido")
		return
	}

	currentID++
	client.ID = strconv.Itoa(currentID)
	clients[client.ID] = client

	respondWithJSON(w, http.StatusCreated, client)
}

// Fiquei em dúvida sobre responser 404 com lista vazia ou 200 com lista vazia
// Segui o padrão de retornar 200 com lista vazia de acordo com boas práticas REST
// https://dev.to/zanfranceschi/conceito-nao-use-http-404-ou-204-para-buscas-sem-resultados-6ki
func getClients(w http.ResponseWriter, r *http.Request) {
	if len(clients) == 0 {
		respondWithJSON(w, http.StatusOK, []Client{})
		return
	}

	var clientList []Client
	for _, client := range clients {
		clientList = append(clientList, client)
	}
	respondWithJSON(w, http.StatusOK, clientList)
}

func getClient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	client, ok := clients[params["id"]]
	if !ok {
		respondWithError(w, http.StatusNotFound, "Cliente não encontrado")
		return
	}
	respondWithJSON(w, http.StatusOK, client)
}

func updateClient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if _, ok := clients[id]; !ok {
		respondWithError(w, http.StatusNotFound, "Cliente não encontrado")
		return
	}

	var clientUpdate Client
	err := json.NewDecoder(r.Body).Decode(&clientUpdate)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Requisição com JSON inválido")
		return
	}

	clientUpdate.ID = id
	clients[id] = clientUpdate

	respondWithJSON(w, http.StatusOK, clientUpdate)
}

// deleteClient exclui um cliente
func deleteClient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if _, ok := clients[id]; !ok {
		respondWithError(w, http.StatusNotFound, "Cliente não encontrado")
		return
	}

	delete(clients, id)
	respondWithJSON(w, http.StatusNoContent, nil)
}

func main() {
	r := mux.NewRouter()

	// Dados iniciais para teste
	clients["1"] = Client{ID: "1", Nome: "John Doe", Nascimento: "01/01/1990", Endereco: "123 Main St", Telefone: "555-5555"}

	// Rotas da API
	r.HandleFunc("/clientes", getClients).Methods("GET")
	r.HandleFunc("/clientes/{id}", getClient).Methods("GET")
	r.HandleFunc("/clientes", createClient).Methods("POST")
	r.HandleFunc("/clientes/{id}", updateClient).Methods("PUT")
	r.HandleFunc("/clientes/{id}", deleteClient).Methods("DELETE")

	log.Println("Servidor iniciado na porta 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
