package errors

var (
	ErrInternalServer     = InternalServer("internal_server", "Terjadi kesalahan sistem")
	ErrDatabaseConnection = InternalServer("database_connection", "Koneksi database gagal")
	ErrUnauthorized       = Unauthorized("unauthorized", "Akses tidak diizinkan")
	ErrForbidden          = Forbidden("forbidden", "Akses ditolak")
)

var (
	ErrCreationFailed  = InternalServer("creation_failed", "Gagal membuat akun")
	ErrUpdateFailed    = InternalServer("update_failed", "failed to update")
	ErrDeletionFailed  = InternalServer("delete_failed", "failed to delete")
	ErrRetrievalFailed = InternalServer("retrieval_failed", "failed to retrieves")
	ErrNotFound        = NotFound("not_found", "not found")
)

var (
	ErrPasswordTooShort = ValidationError("password_too_short", "Password minimal 8 karakter")
	ErrPasswordNumber   = ValidationError("password_no_number", "Password harus mengandung minimal 1 angka")
	ErrPasswordUpper    = ValidationError("password_no_upper", "Password harus mengandung huruf kapital")
	ErrPasswordSpecial  = ValidationError("password_no_special", "Password harus mengandung simbol")
	ErrPasswordNotSame  = ValidationError("password_mismatch", "Password dan konfirmasi password tidak sama")
)

var (
	ErrInvalidCredentials   = Unauthorized("invalid_credentials", "Username atau password salah. Coba lagi!")
	ErrTooManyLoginAttempts = TooManyRequests("too_many_login_attempts", "Terlalu banyak percobaan login. Coba lagi dalam 15 menit")
	ErrUserInactive         = Forbidden("user_inactive", "Akun belum aktif. Silakan lakukan verifikasi email")
	ErrUserNotFound         = NotFound("user_not_found", "Pengguna tidak ditemukan")
	ErrUsernameExists       = Conflict("username_exists", "Username sudah tersedia, silakan gunakan yang lain")
)

var (
	ErrEmailNotRegistered   = ValidationError("email_not_registered", "Email belum terdaftar. Silakan lakukan pendaftaran")
	ErrEmailAlreadyVerified = ValidationError("email_already_verified", "Email sudah diverifikasi")
	ErrEmailExists          = Conflict("email_exists", "Email sudah terdaftar, gunakan alamat email lainnya")
	ErrEmailNotFound        = NotFound("email_not_found", "Email tidak ditemukan")
)

var (
	ErrGenerateToken       = InternalServer("generate_token_failed", "Gagal membuat token")
	ErrInvalidToken        = BadRequest("invalid_token", "Token kadaluwarsa atau sudah digunakan")
	ErrTokenExpired        = Unauthorized("token_expired", "Token sudah kadaluwarsa")
	ErrSaveRefreshToken    = InternalServer("save_refresh_token_failed", "Gagal menyimpan refresh token")
	ErrInvalidRefreshToken = Unauthorized("invalid_refresh_token", "Refresh token tidak valid atau sudah dicabut")
)
