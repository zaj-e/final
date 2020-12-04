package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	puerto_generador_avisar_nuevos_miembros = 1000
	puerto_avisar_nuevos_miembros = 8001
	puerto_recibir_nuevos_miembros = 8002
	puerto_recibir_pokemones = 2000
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

var (
	direccionGenerador string
	yo                 Entrenador
	clubEntrenadores   []Entrenador
	)

func solicitarRegistro(dir string)  {
	con, _ := net.Dial("tcp", dir)
	defer con.Close()
	bMs, _ := json.Marshal(yo)
	fmt.Fprintf(con, string(bMs))
}

func PrimeroEnEscucharNuevoMiembro() {
	host := fmt.Sprintf("%s:%d", yo.direccion, puerto_avisar_nuevos_miembros)
	ln, _ := net.Listen("tcp", host)
	defer ln.Close()
	for {
		con, _ := ln.Accept()
		go AvisarATodosNuevoMiembro(con)
	}

}

//func AvisarNuevoMiembro(entrenador Entrenador, nuevoEntrenador Entrenador) {
//	remoteHost := fmt.Sprintf("%s:%d", entrenador.direccion, puerto_recibir_nuevos_miembros)
//	con, _ := net.Dial("tcp", remoteHost)
//	defer con.Close()
//	bMs, _ := json.Marshal(nuevoEntrenador)
//	fmt.Fprintf(con, string(bMs))
//}

func AvisarGenerador(nuevoEntrenador Entrenador) {
	remotehost := fmt.Sprintf("%s:%d", direccionGenerador, puerto_generador_avisar_nuevos_miembros)
	con, _ := net.Dial("tcp", remotehost)
	defer con.Close()
	bMs, _ := json.Marshal(nuevoEntrenador)
	fmt.Fprintf(con, string(bMs))
}

func AvisarATodosNuevoMiembro(con net.Conn) {
	var nuevoEntrenador Entrenador
	bufferIn := bufio.NewReader(con)
	msgCon, _ := bufferIn.ReadString('\n')
	json.Unmarshal([]byte(msgCon), &nuevoEntrenador)


	go AvisarGenerador(nuevoEntrenador)

	clubEntrenadores = append(clubEntrenadores, nuevoEntrenador)
}


func ManejarAceptarNuevoMiembro(con net.Conn) {
	fmt.Println("Aceptas registrar a este nuevo miembro? : ")
	bIn := bufio.NewReader(os.Stdin)
	respuestaAceptar, _ := bIn.ReadString('\n')
	respuestaAceptar = strings.TrimSpace(respuestaAceptar)

	if (respuestaAceptar == "si"){
		var nuevoEntrenador Entrenador
		bufferIn := bufio.NewReader(con)
		msg, _ := bufferIn.ReadString('\n')
		json.Unmarshal([]byte(msg), nuevoEntrenador)
		clubEntrenadores = append(clubEntrenadores, nuevoEntrenador)
	}
}


func EscucharNuevoMiembro() {

	hostName := fmt.Sprintf("%s:%d", yo.direccion, puerto_recibir_nuevos_miembros)
	ln, _ := net.Listen("tcp", hostName)
	defer ln.Close()

	for {
		con, _ := ln.Accept()
		go ManejarAceptarNuevoMiembro(con)
	}
}


func ManejarCaptura(con net.Conn) {
	var pokemonRecibido Pokemon
	bufferIn := bufio.NewReader(con)
	msg, _ := bufferIn.ReadString('\n')
	json.Unmarshal([]byte(msg), &pokemonRecibido)

	if (pokemonRecibido.tipo == yo.tipoElegido) {
		yo.pokemons = append(yo.pokemons, pokemonRecibido)
	}
}

func RecibirPokemones() {
	hostName := fmt.Sprintf("%s:%d", direccionGenerador, puerto_recibir_pokemones)
	ln, _ := net.Listen("tcp", hostName)
	defer ln.Close()
	for {
		con, _ := ln.Accept()
		go ManejarCaptura(con)
	}
}


func main () {

	fmt.Println("Ingrese direccion del nodo generador: ")
	bIn := bufio.NewReader(os.Stdin)
	direccionGenerador, _ = bIn.ReadString('\n')
	direccionGenerador = strings.TrimSpace(direccionGenerador)

	fmt.Println("Usuario, ingrese la direccion para este entrenador: ")
	bInDos := bufio.NewReader(os.Stdin)
	miDireccion, _ := bInDos.ReadString('\n')
	miDireccion = strings.TrimSpace(miDireccion)
	yo.direccion = miDireccion


	fmt.Println("Escoja el tipo de pokemon a capturar (3, 4, 5 o 6): ")
	bInTres := bufio.NewReader(os.Stdin)
	tipoPokemonElegido , _ := bInTres.ReadString('\n')
	tipoPokemonElegido = strings.TrimSpace(tipoPokemonElegido)
	yo.tipoElegido, _ = strconv.Atoi(tipoPokemonElegido)




	solicitarRegistro(direccionGenerador)

	PrimeroEnEscucharNuevoMiembro()
	EscucharNuevoMiembro()
	RecibirPokemones()
	
}




