package handlers

import (
	"html/template"
	"meu-servidor-go/storage"
	"net/http"
)

func render(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("templates/layout.html", "templates/"+tmpl)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	t.ExecuteTemplate(w, "layout", data)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) { render(w, "index.html", nil) }

func FormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		render(w, "form.html", nil)
		return
	}
	r.ParseForm()
	storage.SalvarMensagem(r.FormValue("nome_usuario"), r.FormValue("mensagem"))
	http.Redirect(w, r, "/ver", http.StatusSeeOther)
}

func ListarHandler(w http.ResponseWriter, r *http.Request) {
	msgs, _ := storage.LerMensagens()
	render(w, "lista.html", msgs)
}

func ExcluirHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id != "" {
		storage.ExcluirMensagem(id)
	}
	http.Redirect(w, r, "/ver", http.StatusSeeOther)
}

func EditarHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Redirect(w, r, "/ver", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		render(w, "editar.html", id)
		return
	}

	r.ParseForm()
	novoNome := r.FormValue("nome")
	novaMsg := r.FormValue("mensagem")
	
	err := storage.AtualizarMensagem(id, novoNome, novaMsg)
	if err != nil {
		http.Error(w, "Erro ao atualizar", 500)
		return
	}

	http.Redirect(w, r, "/ver", http.StatusSeeOther)
}