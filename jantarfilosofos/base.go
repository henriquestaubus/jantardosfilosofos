package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type Garfo chan bool

func filosofo(id int, garfoEsq, garfoDir Garfo, wg *sync.WaitGroup, R int) {
	for i := 0; i < R; i++ {
		// 1. SEÇÃO NÃO-CRÍTICA: Pensar
		fmt.Printf("[Filósofo %d] PENSANDO...\n", id)
		time.Sleep(100 * time.Millisecond)

		// 2. PREPROTOCOL: Tentar pegar os garfos da esquerda e direita
		fmt.Printf("[Filósofo %d] FOME! Tentando garfo ESQUERDO...\n", id)
		<-garfoEsq
		fmt.Printf("[Filósofo %d] Pegou garfo ESQUERDO.\n", id)

		// Pausa intencional para induzir o deadlock
		time.Sleep(10 * time.Millisecond)

		fmt.Printf("[Filósofo %d] Tentando garfo DIREITO...\n", id)
		<-garfoDir
		fmt.Printf("[Filósofo %d] Pegou garfo DIREITO.\n", id)

		// 3. SEÇÃO CRÍTICA: Comer
		fmt.Printf("[Filósofo %d] COMENDO 🍝 (%d/%d)\n", id, i+1, R)
		time.Sleep(150 * time.Millisecond)

		// 4. POSTPROTOCOL: Largar garfos
		garfoEsq <- true
		garfoDir <- true
		fmt.Printf("[Filósofo %d] Largou os garfos.\n", id)
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

	fmt.Printf("=== JANTAR DOS FILÓSOFOS - VERSÃO BASE (N=%d, R=%d) ===\n", N, R)

	var wg sync.WaitGroup
	garfos := make([]Garfo, N)
	for i := 0; i < N; i++ {
		garfos[i] = make(Garfo, 1)
		garfos[i] <- true
	}

	for i := 0; i < N; i++ {
		wg.Add(1)
		garfoEsq := garfos[i]
		garfoDir := garfos[(i+1)%N]
		go filosofo(i, garfoEsq, garfoDir, &wg, R)
	}

	wg.Wait()
	fmt.Println("=== FIM DO JANTAR ===")
}