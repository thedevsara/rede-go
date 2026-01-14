package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Post struct {
	ID               string
	Data             string
	Nome             string
	Mensagem         string
	TotalCurtidas    int
	TotalComentarios int
	Comentarios      []Comentario
}

type Comentario struct {
	ID       int
	Autor    string
	Conteudo string
	Data     string
}

const connStr = "postgres://postgres:sara123@localhost:5432/postgres-go"

func ObterConexao() (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), connStr)
}

// Função centralizada para formatar a data e ajustar o fuso horário (-3h)
func FormatarTempo(t time.Time) string {
	fusoBrasilia := time.FixedZone("BRT", -3*3600)
	horaLocal := t.In(fusoBrasilia)
	return horaLocal.Format("02/01/2006 15:04")
}

// --- LOGICA DE POSTS ---

func SalvarMensagem(nome, mensagem string) error {
	conn, err := ObterConexao()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(),
		"INSERT INTO posts (nome, mensagem, data_hora) VALUES ($1, $2, $3)",
		nome, mensagem, time.Now().UTC()) // Salvamos em UTC para evitar conflitos
	return err
}

func LerMensagens() ([]Post, error) {
	conn, err := ObterConexao()
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	query := `
        SELECT p.id, p.data_hora, p.nome, p.mensagem,
        (SELECT COUNT(*) FROM curtidas WHERE post_id = p.id) as curtidas,
        (SELECT COUNT(*) FROM comentarios WHERE post_id = p.id) as comentarios
        FROM posts p ORDER BY p.id DESC`

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		var id int
		var dataDoBanco time.Time
		err := rows.Scan(&id, &dataDoBanco, &p.Nome, &p.Mensagem, &p.TotalCurtidas, &p.TotalComentarios)
		if err != nil {
			continue
		}

		p.ID = fmt.Sprintf("%d", id)
		p.Data = FormatarTempo(dataDoBanco) // Aplica fuso horário nos posts

		p.Comentarios, _ = LerComentariosDoPost(p.ID)

		posts = append(posts, p)
	}
	return posts, nil
}

// --- CURTIDAS (REGRA: APENAS UMA VEZ) ---

func AlternarCurtida(postID string, usuario string) error {
	conn, err := ObterConexao()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	var existe bool
	err = conn.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM curtidas WHERE post_id = $1 AND usuario_nome = $2)",
		postID, usuario).Scan(&existe)

	if existe {
		_, err = conn.Exec(context.Background(), "DELETE FROM curtidas WHERE post_id = $1 AND usuario_nome = $2", postID, usuario)
	} else {
		_, err = conn.Exec(context.Background(), "INSERT INTO curtidas (post_id, usuario_nome) VALUES ($1, $2)", postID, usuario)
	}
	return err
}

// --- COMENTÁRIOS ---

func AdicionarComentario(postID, autor, conteudo string) error {
	conn, err := ObterConexao()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(),
		"INSERT INTO comentarios (post_id, autor, conteudo, data_criacao) VALUES ($1, $2, $3, $4)",
		postID, autor, conteudo, time.Now().UTC())
	return err
}

func LerComentariosDoPost(postID string) ([]Comentario, error) {
	conn, err := ObterConexao()
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(),
		"SELECT autor, conteudo, data_criacao FROM comentarios WHERE post_id = $1 ORDER BY data_criacao ASC", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []Comentario
	for rows.Next() {
		var c Comentario
		var data time.Time
		if err := rows.Scan(&c.Autor, &c.Conteudo, &data); err == nil {
			c.Data = FormatarTempo(data) // AGORA OS COMENTÁRIOS TAMBÉM USAM O FUSO CORRETO
			lista = append(lista, c)
		}
	}
	return lista, nil
}

// --- USUÁRIOS E SEGURANÇA ---

func CriarUsuario(username, password string) error {
	conn, err := ObterConexao()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err = conn.Exec(context.Background(), "INSERT INTO usuarios (username, password_hash) VALUES ($1, $2)", username, string(hash))
	return err
}

func VerificarUsuario(username, password string) (bool, error) {
	conn, err := ObterConexao()
	if err != nil {
		return false, err
	}
	defer conn.Close(context.Background())
	var hash string
	err = conn.QueryRow(context.Background(), "SELECT password_hash FROM usuarios WHERE username = $1", username).Scan(&hash)
	if err != nil {
		return false, nil
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil, nil
}

// --- MANUTENÇÃO ---

func ExcluirMensagem(id string) error {
	conn, err := ObterConexao()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	_, err = conn.Exec(context.Background(), "DELETE FROM posts WHERE id = $1", id)
	return err
}

func AtualizarMensagem(id, novoNome, novaMsg string) error {
	conn, err := ObterConexao()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	_, err = conn.Exec(context.Background(), "UPDATE posts SET nome = $1, mensagem = $2 WHERE id = $3", novoNome, novaMsg, id)
	return err
}