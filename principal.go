package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

const (
	user      = "postgres"
	password  = "1234"
	baseDados = "sist_distrib_dib"
	host      = "localhost"
	port      = "15432"
)

type Cliente struct {
	Id           int64     `json:"id,omitempty"`
	PrimeiroNome string    `json:"nome,omitempty"`
	Sobrenome    string    `json:"sobrenome,omitempty"`
	Endereco     *Endereco `json:"endereco,omitempty"`
}
type Endereco struct {
	Id         int64  `json:"endereco_id,omitempty"`
	Logradouro string `json:"logradouro,omitempty"`
	Cep        int64  `json:"cep,omitempty"`
	Bairro     string `json:"bairro,omitempty"`
	Cidade     string `json:"cidade,omitempty"`
	UF         string `json:"uf,omitempty"`
}

func conectaNoBancoDeDados() *sql.DB {
	conexao := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, baseDados, host, port)
	db, err := sql.Open("postgres", conexao)
	mensagemErro(err)
	return db
}
func mensagemErro(err error) {
	if err != nil {
		panic(err.Error())
	}
}
func ListarClientes(w http.ResponseWriter, r *http.Request) {
	db := conectaNoBancoDeDados()
	listaCliente, err := db.Query("select c.id, c.primeiro_nome, c.sobrenome, c.endereco_id,  e.logradouro, e.bairro, e.cep, e.cidade, e.uf from tb_cliente  c left join tb_endereco e on c.endereco_id=e.id")
	mensagemErro(err)
	var clientes []Cliente
	for listaCliente.Next() {
		var (
			id, endereco_id, cep                            sql.NullInt64
			nome, sobrenome, logradouro, cidade, bairro, uf sql.NullString
		)
		err := listaCliente.Scan(&id, &nome, &sobrenome, &endereco_id, &logradouro, &bairro, &cep, &cidade, &uf)
		mensagemErro(err)
		cliente := Cliente{}
		cliente.Id = id.Int64
		cliente.PrimeiroNome = nome.String
		cliente.Sobrenome = sobrenome.String
		endereco := Endereco{endereco_id.Int64, logradouro.String, cep.Int64, bairro.String, cidade.String, uf.String}
		cliente.Endereco = &endereco
		clientes = append(clientes, cliente)
	}
	defer db.Close()
	fmt.Println(clientes)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientes)
}

func InsererUmCliente(w http.ResponseWriter, r *http.Request) {
	var cliente Cliente
	json.NewDecoder(r.Body).Decode(&cliente)
	db := conectaNoBancoDeDados()
	var id_endereco sql.NullInt64

	if cliente.Endereco != nil {
		id_endereco = sql.NullInt64{cliente.Endereco.Id, true}
	}
	queryInsert, err := db.Prepare("insert into tb_cliente (id, primeiro_nome, sobrenome, endereco_id) values (nextval('seq_cliente'), $1, $2, $3)")
	mensagemErro(err)
	queryInsert.Exec(cliente.PrimeiroNome, cliente.Sobrenome, id_endereco)
	defer db.Close()
}

func AlterarUmCliente(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	var clienteNovo Cliente

	err := json.NewDecoder(r.Body).Decode(&clienteNovo)
	mensagemErro(err)
	db := conectaNoBancoDeDados()
	alterarCliente, err := db.Prepare("update tb_cliente set primeiro_nome=$1, sobrenome=$2, endereco_id=$3 where id=$4")
	mensagemErro(err)

	var id_endereco sql.NullInt64
	if clienteNovo.Endereco != nil {
		id_endereco = sql.NullInt64{clienteNovo.Endereco.Id, true}
	}
	alterarCliente.Exec(clienteNovo.PrimeiroNome, clienteNovo.Sobrenome, id_endereco, parametros["ID"])
	defer db.Close()

}

func DeletarUmCliente(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	db := conectaNoBancoDeDados()
	deletarCliente, err := db.Prepare("delete from tb_cliente where id=$1")
	mensagemErro(err)
	deletarCliente.Exec(parametros["ID"])
	defer db.Close()

}

func ConsultarPorNomeSobrenome(w http.ResponseWriter, r *http.Request) {
	nomeURL := r.URL.Query().Get("nome")
	sobrenomeURL := r.URL.Query().Get("sobrenome")
	db := conectaNoBancoDeDados()
	consultaNomeSobrenome, err := db.Query("select c.id, c.primeiro_nome, c.sobrenome, c.endereco_id,  e.logradouro, e.bairro, e.cep, e.cidade, e.uf from tb_cliente  c left join tb_endereco e on c.endereco_id=e.id where c.primeiro_nome=$1 and c.sobrenome=$2", nomeURL, sobrenomeURL)
	mensagemErro(err)
	var clientes []Cliente
	for consultaNomeSobrenome.Next() {
		var (
			id, endereco_id, cep                            sql.NullInt64
			nome, sobrenome, logradouro, cidade, bairro, uf sql.NullString
		)
		err := consultaNomeSobrenome.Scan(&id, &nome, &sobrenome, &endereco_id, &logradouro, &bairro, &cep, &cidade, &uf)
		mensagemErro(err)
		cliente := Cliente{}
		cliente.Id = id.Int64
		cliente.PrimeiroNome = nome.String
		cliente.Sobrenome = sobrenome.String
		endereco := Endereco{endereco_id.Int64, logradouro.String, cep.Int64, bairro.String, cidade.String, uf.String}
		cliente.Endereco = &endereco
		clientes = append(clientes, cliente)
	}
	defer db.Close()
	fmt.Println(clientes)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientes)
}

func ConsultaPorUmID(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	idCliente := parametros["ID"]
	db := conectaNoBancoDeDados()
	selecinarUmCliente, err := db.Query("select c.id, c.primeiro_nome, c.sobrenome, c.endereco_id,  e.logradouro, e.bairro, e.cep, e.cidade, e.uf from tb_cliente  c left join tb_endereco e on c.endereco_id=e.id where c.id=$1", idCliente)
	mensagemErro(err)
	cliente := Cliente{}
	for selecinarUmCliente.Next() {
		var (
			id, endereco_id, cep                            sql.NullInt64
			nome, sobrenome, logradouro, cidade, bairro, uf sql.NullString
		)
		err := selecinarUmCliente.Scan(&id, &nome, &sobrenome, &endereco_id, &logradouro, &bairro, &cep, &cidade, &uf)
		mensagemErro(err)
		cliente.Id = id.Int64
		cliente.PrimeiroNome = nome.String
		cliente.Sobrenome = sobrenome.String
		endereco := Endereco{endereco_id.Int64, logradouro.String, cep.Int64, bairro.String, cidade.String, uf.String}
		cliente.Endereco = &endereco
	}
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cliente)
}

func main() {

	rotas := mux.NewRouter()
	rotas.HandleFunc("/clientes", ListarClientes).Methods("GET")
	rotas.HandleFunc("/cliente/{ID}", ConsultaPorUmID).Methods("GET")
	rotas.HandleFunc("/cliente/nome/", ConsultarPorNomeSobrenome).Methods("GET")
	rotas.HandleFunc("/cliente", InsererUmCliente).Methods("POST")
	rotas.HandleFunc("/cliente/{ID}", AlterarUmCliente).Methods("PUT")
	rotas.HandleFunc("/cliente/{ID}", DeletarUmCliente).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8081", rotas))
}
