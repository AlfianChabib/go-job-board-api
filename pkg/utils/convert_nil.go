package utils

// NilIfEmpty mengembalikan nil jika pointer bernilai nil atau menunjuk ke zero value dari tipe T.
// Jika pointer memiliki nilai (bukan zero value), maka akan mengembalikan pointer aslinya.
func NilIfEmpty[T comparable](v *T) *T {
	var zero T
	if v != nil && *v == zero {
		return nil
	}
	return v
}

// NilIfZero adalah alias untuk NilIfEmpty agar lebih semantik saat digunakan untuk angka/boolean/struct.
func NilIfZero[T comparable](v *T) *T {
	return NilIfEmpty(v)
}

// CleanseEmpty mengubah pointer yang ditunjuk menjadi nil jika nilainya adalah zero value.
func CleanseEmpty[T comparable](v **T) {
	var zero T
	if v != nil && *v != nil && **v == zero {
		*v = nil
	}
}

// CleanseEmptyString dipertahankan untuk backward compatibility dengan string pointer.
func CleanseEmptyString(s **string) {
	CleanseEmpty(s)
}
