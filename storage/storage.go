package storage

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Post struct {
	ID       string
	Data     string
	Nome     string
	Mensagem string
}

func SalvarMensagem(nome, mensagem string) error {
	f, err := os.OpenFile("dados.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	id := time.Now().UnixNano()
	dataHora := time.Now().Format("02/01 15:04")
	linha := fmt.Sprintf("%d|%s|%s|%s\n", id, dataHora, nome, mensagem)
	_, err = f.WriteString(linha)
	return err
}

func LerMensagens() ([]Post, error) {
	dados, err := os.ReadFile("dados.txt")
	if err != nil {
		if os.IsNotExist(err) { return []Post{}, nil }
		return nil, err
	}
	linhas := strings.Split(strings.TrimSpace(string(dados)), "\n")
	var posts []Post
	for _, linha := range linhas {
		p := strings.Split(linha, "|")
		if len(p) < 4 { continue }
		posts = append(posts, Post{ID: p[0], Data: p[1], Nome: p[2], Mensagem: p[3]})
	}
	for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
		posts[i], posts[j] = posts[j], posts[i]
	}
	return posts, nil
}

func ExcluirMensagem(id string) error {
	dados, _ := os.ReadFile("dados.txt")
	linhas := strings.Split(strings.TrimSpace(string(dados)), "\n")
	var novasLinhas []string
	for _, linha := range linhas {
		if !strings.HasPrefix(linha, id+"|") {
			novasLinhas = append(novasLinhas, linha)
		}
	}
	return os.WriteFile("dados.txt", []byte(strings.Join(novasLinhas, "\n")+"\n"), 0644)
}


func AtualizarMensagem(id, novoNome, novaMsg string) error {
	dados, err := os.ReadFile("dados.txt")
	if err != nil { return err }

	linhas := strings.Split(strings.TrimSpace(string(dados)), "\n")
	for i, linha := range linhas {
		if strings.HasPrefix(linha, id+"|") {
			partes := strings.Split(linha, "|")
			linhas[i] = fmt.Sprintf("%s|%s|%s|%s", partes[0], partes[1], novoNome, novaMsg)
			break
		}
	}
	return os.WriteFile("dados.txt", []byte(strings.Join(linhas, "\n")+"\n"), 0644)
}