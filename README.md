# 🚀 Go Social - Premium Dark Edition

Um sistema de rede social minimalista e futurista desenvolvido em **Go (Golang)**, focado em alta performance, segurança e uma interface imersiva.

## 📱 Visão Geral
O **Go Social** permite partilhar pensamentos em tempo real, interagir com curtidas únicas e comentar publicações. O design segue a estética *Glassmorphism* e *Deep UI*, proporcionando uma experiência moderna e responsiva.

## ✨ Funcionalidades
- **Autenticação**: Login e Registo com **Bcrypt** e sessões **JWT**.
- **Fuso Horário**: Datas ajustadas automaticamente para o horário de Brasília (-3h).
- **Engajamento**: Curtidas únicas por utilizador e sistema de comentários.
- **UI/UX**: Design Dark mode com efeitos de vidro e navegação flutuante.

## 🛠️ Tecnologias
- **Backend**: Go (Golang)
- **Base de Dados**: PostgreSQL
- **Frontend**: Go Templates & Bootstrap 5

---

## 🚀 Como Executar o Projeto

### 1. Configurar o Banco de Dados 🗄️
Certifica-te de que o **PostgreSQL** está ativo e cria a base de dados. As tabelas devem seguir a estrutura do ficheiro `storage.go`:

| Tabela | Descrição |
| :--- | :--- |
| `usuarios` | Armazena credenciais e hashes de senha. |
| `posts` | Registra as mensagens e metadados. |
| `curtidas` | Controla interações únicas por post/user. |
| `comentarios` | Gerencia as respostas de cada postagem. |

> **Nota:** Verifica a string de conexão em `storage.go` para garantir que o utilizador e a senha coincidem com o teu ambiente local.

### 2. Iniciar o Servidor ⚡
No terminal, dentro da pasta raiz do projeto, executa o comando abaixo para compilar e rodar a aplicação:

```bash
go run main.go
