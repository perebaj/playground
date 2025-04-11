package main

func tmp[V int64 | float64](x V) {
	switch v := (interface{})(x).(type) {
	case int64:
		println("int64", v)
	case float64:
		println("float64", v)
	}
}

func main() {
	tmp(int64(1))
	tmp(1.0)
}
