package varusage

type VarUsage int

const (
  DECLARED VarUsage = iota
  INITIALIZED
  USED
)