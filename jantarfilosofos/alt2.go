// Estratégia 2: Limitação do número de filósofos simultaneamente à mesa (semáforo)
package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type Garfo chan bool

type Estatisticas struct {
	refeicoesFeitas  int
	tempoTotalEspera time.Duration
}

func filosofo(id int, garfoEsq, garfoDir Garfo, mesa chan struct{}, wg *sync.WaitGroup, R int, stats *Estatisticas, mu *sync.Mutex) {
	for i := 0; i < R; i++ {
		// 1. SEÇÃO NÃO-CRÍTICA: Pensar
		fmt.Printf("[Filósofo %d] PENSANDO...\n", id)
		time.Sleep(100 * time.Millisecond)

		inicioEspera := time.Now()
		fmt.Printf("[Filósofo %d] com FOME. Aguardando vaga na mesa...\n", id)

		// 2. CONTROLE DE ACESSO À MESA (SEMÁFORO - MÁXIMO N-1 FILÓSOFOS)
		mesa <- struct{}{}
		fmt.Printf("[Filósofo %d] Sentou à mesa.\n", id)

		// 3. PREPROTOCOL: Pegar garfos
		<-garfoEsq
		<-garfoDir

		tempoEspera := time.Since(inicioEspera)

		// Atualização das métricas protegida por Mutex
		mu.Lock()
		stats.refeicoesFeitas++
		stats.tempoTotalEspera += tempoEspera
		mu.Unlock()

		// 4. SEÇÃO CRÍTICA: Comer
		fmt.Printf("[Filósofo %d] COMENDO 🍝 (Refeição %d/%d) (Esperou: %v)\n", id, i+1, R, tempoEspera)
		time.Sleep(150 * time.Millisecond)

		// 5. POSTPROTOCOL: Largar garfos e sair da mesa
		garfoEsq <- true
		garfoDir <- true

		<-mesa
		fmt.Printf("[Filósofo %d] Levantou da mesa e voltou a pensar.\n", id)
	}
	wg.Done()
}

func main() {
	N := 5
	R := 3

	if len(os.Args) > 1 {
		if val, err := strconv.Atoi(os.Args[1]); err == nil && val > 0 {
			N = val
		}
	}
	if len(os.Args) > 2 {
		if val, err := strconv.Atoi(os.Args[2]); err == nil && val > 0 {
			R = val
		}
	}

	fmt.Printf("=== JANTAR (ESTRATÉGIA 2: LIMITAÇÃO DE MESA - N=%d, R=%d) ===\n", N, R)

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Semáforo limitando a N-1 assentos na mesa
	mesa := make(chan struct{}, N-1)

	garfos := make([]Garfo, N)
	for i := 0; i < N; i++ {
		garfos[i] = make(Garfo, 1)
		garfos[i] <- true
	}

	estatisticas := make([]Estatisticas, N)

	for i := 0; i < N; i++ {
		wg.Add(1)
		garfoEsq := garfos[i]
		garfoDir := garfos[(i+1)%N]
		go filosofo(i, garfoEsq, garfoDir, mesa, &wg, R, &estatisticas[i], &mu)
	}

	wg.Wait()

	fmt.Println("\n=======================================================")
	fmt.Println("=== RELATÓRIO DE EXECUÇÃO E MÉTRICAS DE FAIRNESS ===")
	fmt.Println("=======================================================")
	for i := 0; i < N; i++ {
		tempoMedio := estatisticas[i].tempoTotalEspera / time.Duration(R)
		fmt.Printf("Filósofo %d: %d refeições concluídas | Tempo Médio de Espera: %v\n",
			i, estatisticas[i].refeicoesFeitas, tempoMedio)
	}
	fmt.Println("=======================================================")
	fmt.Println("Encerramento limpo: Nenhuma goroutine bloqueada.")
}