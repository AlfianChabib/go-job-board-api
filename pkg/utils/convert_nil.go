package utils

func NilIfEmpty(s *string) *string {
	if s != nil && *s == "" {
		return nil
	}
	return s
}

func CleanseEmptyString(s **string) {
	// Pengecekan aman untuk mencegah panic (nil pointer dereference)
	if s != nil && *s != nil && **s == "" {
		*s = nil // Timpa pointer asli menjadi nil
	}
}
