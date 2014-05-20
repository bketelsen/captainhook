package types

type Orchestration struct {
	Scripts []Script `json:"scripts"`
}

type Script struct {
	Command string `json:"command"`
}
