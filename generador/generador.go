package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"
)

const (
	puerto_generador_avisar_nuevos_miembros = 1000
	puerto_recibir_nuevos_miembros = 8002
	)

var (
	clubEntrenadores   []Entrenador
	miDireccion string
)

type Pokemon struct {
	id int
	//nombre string
	tipo int
}

type Entrenador struct {
	direccion string
	pokemons []Pokemon
	tipoElegido int
}

func main()  {

	go GenerarPokemones()
	EscucharNuevoMiembroComoGenerador()
}

func GenerarPokemones() {
	for {
		time.Sleep(5 * time.Second)
		var nuevoPokemon Pokemon
		nuevoPokemon.id = rand.Intn(10000000)
		nuevoPokemon.tipo = rand.Intn(7 - 3) + 3 // entre 3 y 6

		indiceEntrenadorAlQueSeEnviaElNuevoPokemon := rand.Intn(len(clubEntrenadores))


	}
}


func AvisarATodosNuevoMiembro(con net.Conn) {
	var nuevoEntrenador Entrenador
	bufferIn := bufio.NewReader(con)
	msgCon, _ := bufferIn.ReadString('\n')
	json.Unmarshal([]byte(msgCon), &nuevoEntrenador)


	for _, entrenador := range clubEntrenadores {
		go AvisarNuevoMiembro(entrenador, nuevoEntrenador)
	}

	clubEntrenadores = append(clubEntrenadores, nuevoEntrenador)
}


func AvisarNuevoMiembro(entrenador Entrenador, nuevoEntrenador Entrenador) {
	remoteHost := fmt.Sprintf("%s:%d", entrenador.direccion, puerto_recibir_nuevos_miembros)
	con, _ := net.Dial("tcp", remoteHost)
	defer con.Close()
	bMs, _ := json.Marshal(nuevoEntrenador)
	fmt.Fprintf(con, string(bMs))
}


func EscucharNuevoMiembroComoGenerador() {
	hostName := fmt.Sprintf("%s:%d", miDireccion, puerto_recibir_nuevos_miembros)

	ln, _ := net.Listen("tcp", hostName)
	defer ln.Close()

	for {
		con, _ := ln.Accept()
		go AvisarATodosNuevoMiembro(con)
	}

}