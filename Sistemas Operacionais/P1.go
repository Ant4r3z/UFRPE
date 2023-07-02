/*DESCRIPTION

1. Cinco threads concorrendo simultaneamente a dois recursos 
compartilhados (ex.: variáveis globais, buffers, etc.), as 
quais somente uma thread por vez pode acessar cada recurso 
compartilhado (pode utilizar qualquer técnica de exclusão mútua). 
Nesse caso, deve ser demonstrado (com logs, prints, graficamente, etc.) 
que as condições de corrida existem e que a exclusão mútua de fato 
ocorre. Como sugrscão, considerar a permuta de uma thread para outra 
a cada 3 segundos (nesse intervalo, uma thread pode rscar "consumindo" 
o recurso enquanto as demais aguardam, seja de forma bloqueada ou como 
espera ocupada), como no exemplo do produtor/consumidor (porém aqui 
com dois buffers).

*/

//////////////////////////////////////////////////////
/*READ ME
	Para evitar que um mesmo Thread consiga useResource
	duas vezes, foi criada uma rscrutura 'Thread' com 
	um booleano para verifithread se ele já rscá resource1ionado.
	Portanto, os métodos 'entrar' e 'sair' recebem um 
	objeto 'Thread' como parâmetro, ao invés do id.
*/

package main

import (
	"fmt"
	"time"
)

type Resource struct { //rscrutura Resource
	id string
	buffer chan string
	TOTAL, free int 
	lock chan int
}

type Thread struct {
	id string
	busy bool
}

func NewResource(id string, total int) Resource {
	var lock chan int = make(chan int, 1)
	var buffer chan string = make(chan string, total)
	r := Resource{id, buffer, total, total, lock}
	return r
}

func NewThread(id string) Thread {
	t := Thread{id, false}
	return t
}

func runResource(rsc *Resource, thread Thread){
	for {
		if thread.busy == false {
			useResource(rsc, thread)
		} else {
			time.Sleep(155*time.Millisecond)
		}
	}
}

func useResource(rsc *Resource, thread Thread){
	rsc.lock <- 0 //uso do lock para bloquear a entrada para outras Threads
	if rsc.free > 0 { //verificacao de buffer livre
		thread.busy = true
		rsc.free-- //diminuicao da quantidade de slots livres no buffer
		rsc.buffer <- thread.id
		fmt.Print(thread.id, " usando ", rsc.id, "\n")
		<- rsc.lock
		sair(rsc, thread)
	} else {
		<- rsc.lock
		time.Sleep(55*time.Millisecond)
	}
}

func sair(rsc *Resource, thread Thread){
	wait()
	rsc.lock <- 1 //uso do lock para bloquear a saída para outras Threads
	fmt.Print(<- rsc.buffer, " liberou ", rsc.id, "\n")
	rsc.free++ //aumento da quantidade de slots livres no buffer
	thread.busy = false
	<- rsc.lock
}

func wait(){ //funcao para esperar 3 segundos antes de liberar
	time.Sleep(3*time.Second)
}


func main() {

	var resource1 = NewResource("Resource1", 2)
	var resource2 = NewResource("Resource2", 2)
	var t1 = NewThread("Thread 1")
	var t2 = NewThread("Thread 2")
	var t3 = NewThread("Thread 3")
	var t4 = NewThread("Thread 4")
	var t5 = NewThread("Thread 5")

	// 

	go runResource(&resource1, t1)
	go runResource(&resource1, t2)
	go runResource(&resource1, t3)
	go runResource(&resource1, t4)
	go runResource(&resource1, t5)
	
	go runResource(&resource2, t1)
	go runResource(&resource2, t2)
	go runResource(&resource2, t3)
	go runResource(&resource2, t4)
	go runResource(&resource2, t5)

	time.Sleep(time.Second * 10)
}
