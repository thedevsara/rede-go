package main

import (
	"fmt"
	"log"
	"meu-servidor-go/handlers"
	"net/http"
)

func main() {
//rotas dos handlers
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/form", handlers.FormHandler)
	http.HandleFunc("/ver", handlers.ListarHandler)
	http.HandleFunc("/excluir", handlers.ExcluirHandler)
	http.HandleFunc("/editar", handlers.EditarHandler)

	fmt.Println("Servidor rodando em http://localhost:8080")
	
	// iniciando o servidor na porta 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}