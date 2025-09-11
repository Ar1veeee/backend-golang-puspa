package errors

var (
	ErrInternalServer     = InternalServer("internal_server", "Terjadi kesalahan sistem")
	ErrDatabaseConnection = InternalServer("database_connection", "Koneksi database gagal")
	ErrTransactionFailed  = InternalServer("transaction_failed", "Transaksi database gagal")
	ErrUniqueViolation    = Conflict("unique_violation", "Data sudah ada")
)

var (
	ErrBadRequest   = BadRequest("bad_request", "Permintaan tidak valid")
	ErrUnauthorized = Unauthorized("unauthorized", "Akses tidak diizinkan")
	ErrForbidden    = Forbidden("forbidden", "Akses ditolak")
	ErrNotFound     = NotFound("not_found", "Data tidak ditemukan")
	ErrConflict     = Conflict("conflict", "Konflik data")
)

var (
	ErrInvalidToken   = BadRequest("invalid_token", "Token tidak valid")
	ErrInvalidInput   = BadRequest("invalid_input", "Input tidak valid")
	ErrInvalidUserID  = BadRequest("invalid_user_id", "Format ID pengguna tidak valid")
	ErrUserIDRequired = BadRequest("user_id_required", "ID pengguna diperlukan")
)

var (
	ErrUserNotFound   = NotFound("user_not_found", "Pengguna tidak ditemukan")
	ErrEmailExists    = Conflict("email_exists", "Email sudah terdaftar, gunakan alamat email lainnya")
	ErrEmailNotFound  = NotFound("email_not_found", "Email tidak ditemukan")
	ErrUsernameExists = Conflict("username_exists", "Username sudah tersedia, silakan gunakan yang lain")
)

var (
	ErrPasswordTooShort = ValidationError("password_too_short", "Password minimal 8 karakter")
	ErrPasswordNumber   = ValidationError("password_no_number", "Password harus mengandung minimal 1 angka")
	ErrPasswordUpper    = ValidationError("password_no_upper", "Password harus mengandung huruf kapital")
	ErrPasswordSpecial  = ValidationError("password_no_special", "Password harus mengandung simbol")
	ErrPasswordNotSame  = ValidationError("password_mismatch", "Password dan konfirmasi password tidak sama")
)
