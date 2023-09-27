package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Process struct {
	ID           int
	Channels     []chan int
	Accounts     [3]int
	LocalState   []int
	TotalBalance int
}

func NewProcess(id int, numProcesses int) *Process {
	channels := make([]chan int, numProcesses)
	localState := make([]int, numProcesses)
	for i := 0; i < numProcesses; i++ {
		channels[i] = make(chan int, 100) 
	}
	return &Process{
		ID:           id,
		Channels:     channels,
		LocalState:   localState,
		TotalBalance: 3 * 1000000, 
	}
}

func (p *Process) transaction(targetID int) {
	amount := rand.Intn(1000) 
	p.Accounts[targetID] += amount
	p.TotalBalance += amount
	fmt.Printf("Process %d: Transferred %d to Process %d\n", p.ID, amount, targetID)
}

func (p *Process) initiateSnapshot(round int) {
	p.LocalState[p.ID] = p.TotalBalance
	for i, ch := range p.Channels {
		if i != p.ID {
			ch <- p.ID 
		}
	}

	time.Sleep(time.Millisecond * 100)

	for i, ch := range p.Channels {
		if i != p.ID {
			<-ch 
			p.LocalState[i] = p.TotalBalance
		}
	}

	fmt.Printf("Snapshot %d for Process %d: Total Balance = %d\n", round, p.ID, p.LocalState)
}

func main() {
	var wg sync.WaitGroup

	numProcesses := 3
	processes := make([]*Process, numProcesses)

	for i := 0; i < numProcesses; i++ {
		processes[i] = NewProcess(i, numProcesses)
	}

	totalBalance := 0
	for _, p := range processes {
		totalBalance += p.TotalBalance
	}
	fmt.Printf("Initial Total Balance: %d\n", totalBalance)

	rand.Seed(time.Now().UnixNano())

	for _, p := range processes {
		wg.Add(1)
		go func(proc *Process) {
			defer wg.Done()
			for i := 0; i < numProcesses*15; i++ {
				targetID := rand.Intn(numProcesses)
				proc.transaction(targetID)
				if i > 0 && i%(numProcesses) == 0 {
					proc.initiateSnapshot(i / numProcesses)
				}
				time.Sleep(time.Millisecond * 10)
			}
		}(p)
	}

	wg.Wait()
}
