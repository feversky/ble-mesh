package def

const (
	ErrorInvalidKey = ErrorType(iota)
	Invalid
	NotFound
)

type ErrorType uint

var errorDict = map[ErrorType]string{
	ErrorInvalidKey: "",
}

func (et ErrorType) Error() string {
	return errorDict[et]
}
