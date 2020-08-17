package main

type counter struct {
	i  chan int
	cc chan int
}

func newCounter() *counter {
	counterChan := make(chan int)
	closeChan := make(chan int)

	ct := &counter{
		i:  counterChan,
		cc: closeChan,
	}

	go func(c *counter) {
		var count = 0
		for {
			select {
			case c.i <- count:
				count++
				if count < 0 {
					count = 0
				}
			case <-c.cc:
				close(c.i)
				close(c.cc)
				return
			}
		}
	}(ct)

	return ct
}

func (ct *counter) Close() {
	ct.cc <- 0
}

func (ct *counter) GetCount() int {
	i := <-ct.i
	return i
}
