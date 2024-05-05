package logger

type Field struct {
	Key   string
	Value any
}

func String(key string, val any) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}

func Error(val error) Field {
	return Field{
		Value: val,
	}
}
