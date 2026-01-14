package handlers

import (
	"html/template"
	"meu-servidor-go/storage"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Chave secreta para assinar o token
var jwtKey = []byte("DZVKSrd1JF5itrRrDOljQgcovfBhTrWd7bwgasYF1P5qAoukKX2yaBRdxKC2UTs9")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func gerarToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func obterUsuarioDoToken(r *http.Request) string {
	cookie, err := r.Cookie("token_jwt")
	if err != nil {
		return ""
	}
	tokenStr := cookie.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !tkn.Valid {
		return ""
	}
	return claims.Username
}

func render(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("templates/layout.html", "templates/"+tmpl)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	t.ExecuteTemplate(w, "layout", data)
}

func ValidarJWT(proximo http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if obterUsuarioDoToken(r) == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		proximo.ServeHTTP(w, r)
	}
}

// --- HANDLERS EXISTENTES ---

func HomeHandler(w http.ResponseWriter, r *http.Request) { render(w, "index.html", nil) }

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		render(w, "login.html", nil)
		return
	}
	r.ParseForm()
	user := r.FormValue("username")
	pass := r.FormValue("password")
	valido, _ := storage.VerificarUsuario(user, pass)
	if valido {
		tokenString, _ := gerarToken(user)
		http.SetCookie(w, &http.Cookie{
			Name: "token_jwt", Value: tokenString, Path: "/",
			Expires: time.Now().Add(24 * time.Hour), HttpOnly: true,
		})
		http.Redirect(w, r, "/ver", http.StatusSeeOther)
		return
	}
	render(w, "login.html", "Usuário ou senha inválidos")
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "token_jwt", Value: "", Path: "/", Expires: time.Unix(0, 0)})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CadastroHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		render(w, "cadastro.html", nil)
		return
	}
	r.ParseForm()
	user, pass := r.FormValue("username"), r.FormValue("password")
	if err := storage.CriarUsuario(user, pass); err != nil {
		render(w, "cadastro.html", "Erro ao criar usuário")
		return
	}
	tokenString, _ := gerarToken(user)
	http.SetCookie(w, &http.Cookie{
		Name: "token_jwt", Value: tokenString, Path: "/",
		Expires: time.Now().Add(24 * time.Hour), HttpOnly: true,
	})
	http.Redirect(w, r, "/ver", http.StatusSeeOther)
}

func FormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		render(w, "form.html", nil)
		return
	}
	storage.SalvarMensagem(obterUsuarioDoToken(r), r.FormValue("mensagem"))
	http.Redirect(w, r, "/ver", http.StatusSeeOther)
}

func ListarHandler(w http.ResponseWriter, r *http.Request) {
	usuarioAtual := obterUsuarioDoToken(r)
	msgs, _ := storage.LerMensagens()
	render(w, "lista.html", map[string]interface{}{
		"Mensagens": msgs,
		"Logado":    usuarioAtual != "",
		"Usuario":   usuarioAtual,
	})
}

// --- HANDLERS DE INTERAÇÃO (SEM RETWEET) ---

func CurtirHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("id")
	usuario := obterUsuarioDoToken(r)
	if postID != "" && usuario != "" {
		storage.AlternarCurtida(postID, usuario)
	}
	// O Referer garante que o usuário não perca o lugar no scroll da página
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func ComentarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		postID := r.FormValue("post_id")
		conteudo := r.FormValue("conteudo")
		autor := obterUsuarioDoToken(r)
		if postID != "" && conteudo != "" {
			storage.AdicionarComentario(postID, autor, conteudo)
		}
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

// --- MANUTENÇÃO ---

func ExcluirHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id != "" {
		storage.ExcluirMensagem(id)
	}
	http.Redirect(w, r, "/ver", http.StatusSeeOther)
}

func EditarHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if r.Method == "GET" {
		render(w, "editar.html", id)
		return
	}
	storage.AtualizarMensagem(id, r.FormValue("nome"), r.FormValue("mensagem"))
	http.Redirect(w, r, "/ver", http.StatusSeeOther)
}