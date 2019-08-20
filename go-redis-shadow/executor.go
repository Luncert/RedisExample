package main

type Executor interface {
	Execute(input byte[], s *Storage) (output byte[])
}
