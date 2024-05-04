package logger

type Field struct {
	Key string
	Val any
}

func String(key string, val any) Field {
	return Field{
		Key: key,
		Val: val,
	}
}

func Error(val error) Field {
	return Field{
		Val: val,
	}
}
