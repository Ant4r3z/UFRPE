/*DESCRIPTION

2. Um "mini" simulador (pode se basear nessa ferramenta: 
https://sourceforge.net/projects/oscsimulator/) de escalonamento 
preemptivo de processos, onde seja possível um usuário 
(não precisa de interface gráfica, pode ser linha de comando):

-   Criar processos indicando: ID, Nome, prioridade, processo 
	I/O bound ou CPU/bound, tempo de CPU total (ex.: em unidades 
	inteiras de tempo, por exemplo, 1 a 10 ms). A cada criação, 
	o processo deve ser inserido na fila de "pronto" para ser 
	escalonado conforme algoritmo de escalonamento;

-   Escolher uma de duas opções de algoritmo de escalonamento 
	implementadas (se em dupla escolher uma por integrante);

-   Selecionar o tempo de quantum da preempção (ex.: em unidades 
	inteiras de tempo, por exemplo, 1 a 10 ms)

-   Mostrar a lista de processos na fila de "prontos" dinamicamente 
	(atualizar conforme escalonamento);

-   Iniciar a execução e escalonamento de processos, mostrando 
	(com logs, prints, graficamente, etc.) ao usuário qual processo 
	está ativo na CPU (por quanto tempo), a preempção do processo e 
	quais estão aguardando, indicando sempre a ordem de execução dos 
	algoritmos.

-   Ao final da execução, indicar o tempo de turnaround de cada 
	processo e o tempo médio de espera de todos os processos.

*/

package main

import (
	"fmt"
	"time"
)

type Process struct {
	ID int
	name string
	priority int
	ioBound bool
	cpuBound bool
	cpuTime int
	turnaround int
	waitingTime int
}

// Funções principais

func CreateProcess(pronto []Process,
				ID int, 
				name string, 
				priority int, 
				ioBound bool, 
				cpuBound bool, 
				cpuTime int) {

	p := Process{ID, name, priority, ioBound, cpuBound, cpuTime, 0, 0}
	pronto = append(pronto, p)
	fmt.Printf("Processo %s criado\n", name)
}

func runRoundRobin(pronto []Process, finalized []Process, quantum int) {

	/*
	É o tipo de escalonamento preemptivo mais simples e consiste 
	em repartir uniformemente o tempo da CPU entre todos os 
	processos prontos para a execução. Os processos são organizados 
	numa fila circular, alocando-se a cada um uma fatia de tempo da CPU, 
	igual a um número inteiro de quantum. Caso um processo não termine 
	dentro de sua fatia de tempo, ele é colocado no fim da fila e uma 
	nova fatia de tempo é alocada para o processo no começo da fila.
	*/
	
	var lock chan int = make(chan int, 1)
	
	n := len(pronto)
	remainingTime := make([]int, n)

	for i := 0; i < n; i++ { // Copia o tempo de execução dos processos em um array
		remainingTime[i] = pronto[i].cpuTime
	}

	t := 0 // Tempo atual
	done := false // Variável que determina o fim de todos os processos
	for !done {
		done = true
		for i := 0; i < n; i++ { // Itera pelos processos

			if remainingTime[i] > 0 {
				done = false

				lock <- 0
				if remainingTime[i] > quantum { // Executa durante o quantum
					fmt.Printf("Processo %s executando entre %d e %d\n", pronto[i].name, t, t+quantum)
					remainingTime[i] -= quantum
					t += quantum
					<- lock
					time.Sleep(time.Millisecond * time.Duration(quantum) * 10)
				} else { // Executa o tempo restante
					fmt.Printf("Processo %s executando entre %d to %d\n", pronto[i].name, t, t+remainingTime[i])
					t += remainingTime[i]
					pronto[i].waitingTime = t - pronto[i].cpuTime
					pronto[i].turnaround = t
					finalized = append(finalized, pronto[i])
					remainingTime[i] = 0
					<- lock
					time.Sleep(time.Millisecond * time.Duration(remainingTime[i]) * 10)
				}

				if isRRDone(remainingTime) { // Checa se todos os processos foram terminados
					done = true
					break
				}
			}
		}
	}
}

func runPriority(pronto []Process, finalized []Process, quantum int) {

	/*
	No escalonamento por prioridades, a cada tarefa é associada uma 
	prioridade (número inteiro) usada para escolher a próxima tarefa 
	a receber o processador, a cada troca de contexto.
	*/
	
	// Tier de prioridades
	var t4 []Process
	var t3 []Process
	var t2 []Process
	var t1 []Process

	// Preencher filas
	for i := range pronto {
		switch pronto[i].priority {
		case 4:
			t4 = append(t4, pronto[i])
		case 3:
			t3 = append(t3, pronto[i])
		case 2:
			t2 = append(t2, pronto[i])
		case 1:
			t1 = append(t1, pronto[i])
		default:
			fmt.Print("Prioridade desconhecida")
		}
	}

	var allProntos = append(t4, t3...)
	allProntos = append(allProntos, t2...)
	allProntos = append(allProntos, t1...)

	for i := range allProntos {
		printProcess(allProntos[i])
	}

	runRoundRobin(t4, finalized, quantum)
	runRoundRobin(t3, finalized, quantum)
	runRoundRobin(t2, finalized, quantum)
	runRoundRobin(t1, finalized, quantum)

}

// Utilitários

func isRRDone(remainingTime []int) bool {
	for _, time := range remainingTime {
		if time > 0 {
			return false
		}
	}
	return true
}

func printProcess(process Process) {
	fmt.Printf("%d	%s	%d	%b	%b	%d	%d	%d", 
				process.ID, 
				process.name, 
				process.priority, 
				process.ioBound, 
				process.cpuBound, 
				process.cpuTime, 
				process.turnaround, 
				process.waitingTime)
}

// Main

func main() {

	var pronto []Process
	var finalized []Process

	var quantum int = 3
	var escAlgo int = 1

	out := false
	for !out {
		fmt.Print("1 - Criar processo\n2 - Escolher algoritmo de escalonamento\n3 - Selecionar quantum\n4 - Iniciar execução\n")

		var option int
		fmt.Scan(&option)

		switch option {
		case 1:
			fmt.Print("Digite os seguintes valores entre espaços: ID, Nome, Prioridade, I/O Bound, CPU Bound, CPU Time\n")
			var ID, prioridade, cpuTime int
			var nome string
			var ioBound, cpuBound bool
			fmt.Scan(&ID, &nome, &prioridade, &ioBound, &cpuBound, &cpuTime)
			CreateProcess(pronto, ID, nome, prioridade, ioBound, cpuBound, cpuTime)

		case 2:
			fmt.Print("1 - Round Robin\n2 - Prioridade\n")
			fmt.Scan(&escAlgo)

		case 3:
			fmt.Print("Quantum: ")
			fmt.Scan(&quantum)

		case 4:
			if (escAlgo == 1) {runRoundRobin(pronto, finalized, quantum)}
			if (escAlgo == 2) {runPriority(pronto, finalized, quantum)} 

			for i := range finalized {
				printProcess(finalized[i])
			}
			
			time.Sleep(time.Second * 15)
		default:
			fmt.Print("Opção inexistente.")
		}
	}
}