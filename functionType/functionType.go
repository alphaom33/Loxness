package functiontype

type FunctionType int

const (
  NONE FunctionType = iota
  FUNCTION
  INITIALIZER
  METHOD
)