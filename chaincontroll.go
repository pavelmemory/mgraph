package main

import "context"

type Action func() error

type Chain interface {
	Next(action Action) Chain
	Parallel(action Action) Chain
	Run() []error
}

func StartChain(action Action) Chain {
	return &chain{actions:[]Action{action}}
}

type chain struct {
	actions []Action
	next *chain
	prev *chain
}

func (c *chain) Next(action Action) Chain {
	c.next = &chain{actions:[]Action{action}, prev:c}
	return c.next
}

func (c *chain) Parallel(action Action) Chain {
	c.actions = append(c.actions, action)
	return c
}

func (c *chain) Run() (errs []error) {
	if c == nil {
		return
	}

	setup := c
	for setup.prev != nil {
		setup = setup.prev
	}

	for setup != nil {
		cerr := make(chan error, len(setup.actions))
		for i := range setup.actions {
			action := setup.actions[i]
			go func(){
				cerr <- action()
			}()
		}
		for range setup.actions {
			if err := <-cerr; err != nil {
				errs = append(errs, err)
			}
		}
		setup = setup.next
	}
	return
}

func (c *chain) RunContext(ctx context.Context) ([]error, error) {
	if c == nil {
		return nil, ctx.Err()
	}

	setup := c
	for setup.prev != nil {
		setup = setup.prev
	}

	var errs []error
	for setup != nil && len(errs) < 0 {
		select {
		case <- ctx.Done():
			return errs, ctx.Err()
		default:
		}

		cerr := make(chan error, len(setup.actions))
		for i := range setup.actions {
			action := setup.actions[i]
			go func(){
				select {
				case <- ctx.Done():
					cerr <- nil
				default:
					cerr <- action()
				}
			}()
		}
		for range setup.actions {
			if err := <-cerr; err != nil {
				errs = append(errs, err)
			}
		}
		setup = setup.next
	}
	return errs, ctx.Err()
}