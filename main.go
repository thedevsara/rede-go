package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"meu-servidor-go/handlers"
	"meu-servidor-go/storage"
)

func main() {
	// 1. Validar a conexão com o banco
	conn, err := storage.ObterConexao()
	if err != nil {
		log.Fatalf("❌ Erro fatal: Não foi possível conectar ao banco: %v", err)
	}
	conn.Close(context.Background())
	fmt.Println("✅ Conexão com PostgreSQL validada!")

	// 2. Definição das Rotas
	
	// Rotas Abertas
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/cadastro", handlers.CadastroHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	
	// Servidor de arquivos estáticos
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// ROTAS PROTEGIDAS POR JWT
	http.HandleFunc("/ver", handlers.ValidarJWT(handlers.ListarHandler))
	http.HandleFunc("/form", handlers.ValidarJWT(handlers.FormHandler))
	http.HandleFunc("/excluir", handlers.ValidarJWT(handlers.ExcluirHandler))
	http.HandleFunc("/editar", handlers.ValidarJWT(handlers.EditarHandler))
	
	// NOVAS ROTAS DE ENGAJAMENTO (Focadas em Curtidas e Comentários)
	http.HandleFunc("/curtir", handlers.ValidarJWT(handlers.CurtirHandler))
	http.HandleFunc("/comentar", handlers.ValidarJWT(handlers.ComentarHandler))
	
	// 3. Inicialização do Servidor
	fmt.Println("🚀 Servidor em http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Erro ao iniciar servidor: ", err)
	}
}