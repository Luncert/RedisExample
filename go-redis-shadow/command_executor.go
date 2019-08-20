package main

//Commands
const {
	DeleteKey = iota
	SetString
	GetString
}

type CommandExecutor struct {
}

func Execute(input byte[], s *Storage) (output byte[]) {
	
}