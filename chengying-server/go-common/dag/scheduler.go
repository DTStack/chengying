package dag

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"sync"
)

type Callback func(Node) error

func Execute(g *Graph, cb Callback) error {
	s := &scheduler{
		waitGroups: make(map[Node]*sync.WaitGroup, 0),
	}
	for _, n := range g.Nodes {
		s.Init(n)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eg, errCtx := errgroup.WithContext(ctx)
	for _, n := range g.Nodes {
		n := n
		errCtx := errCtx
		eg.Go(func() error {
			err := s.routine(n, g.DirectDependees(n), cb, errCtx)
			if err != nil {
				fmt.Printf("node: %d, err: %v\n", int(n), err)
				cancel()
				return err
			}
			return err
		})
	}
	return eg.Wait()
}

func CheckGoroutineError(errContext context.Context) error {
	select {
	case <-errContext.Done():
		return errContext.Err()
	default:
		return nil
	}
}

type scheduler struct {
	waitGroups map[Node]*sync.WaitGroup
}

func (s *scheduler) Init(n Node) {
	s.waitGroups[n] = &sync.WaitGroup{}
	s.waitGroups[n].Add(1)
}

func (s *scheduler) WaitForComplete(n Node) {
	s.waitGroups[n].Wait()
}

func (s *scheduler) routine(n Node, dependency []Node, cb Callback, errCtx context.Context) error {
	defer s.waitGroups[n].Done()
	for _, dep := range dependency {
		s.WaitForComplete(dep)
	}
	err := CheckGoroutineError(errCtx)
	if err != nil {
		return err
	}
	err = cb(n)
	return err
}
