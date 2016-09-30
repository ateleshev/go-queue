package queue

import (
	"sync"
	//	"log"
)

type Pools struct {
	sync.Mutex

	index int
	data  []*Pool
}

func (this *Pools) Len() int { // {{{
	return len(this.data)
} // }}}

func (this *Pools) Less(i, j int) bool { // {{{
	return this.data[i].Pending() < this.data[j].Pending()
} // }}}

func (this *Pools) Swap(i, j int) { // {{{
	defer this.Unlock()
	this.Lock()

	this.data[i], this.data[j] = this.data[j], this.data[i]
} // }}}

func (this *Pools) Push(pool *Pool) { // {{{
	defer this.Unlock()
	this.Lock()

	this.data = append(this.data, pool)
} // }}}

func (this *Pools) Last() int { // {{{
	return this.Len() - 1
} // }}}

func (this *Pools) Pop() *Pool { // {{{
	defer this.Unlock()
	this.Lock()

	n := this.Last()
	pool := this.data[n]
	this.data = this.data[:n]

	return pool
} // }}}

func (this *Pools) Index() int { // {{{
	defer func() {
		// Last() = len - 1
		if this.index < this.Last() {
			this.index++
		} else {
			this.index = 0
		}
		this.Unlock()
	}()
	this.Lock()

	return this.index
} // }}}

func (this *Pools) Get() *Pool { // {{{
	return this.data[this.Index()]
} // }}}

func (this *Pools) Dispatch(job JobInterface) { // {{{
	this.Get().Dispatch(job)
} // }}}
