package seeds

func strOrNil(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
