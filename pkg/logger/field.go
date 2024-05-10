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

func Int64(key string, val int64) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}

func Error(val error) Field {
	return Field{
		Key:   "error",
		Value: val,
	}
}
