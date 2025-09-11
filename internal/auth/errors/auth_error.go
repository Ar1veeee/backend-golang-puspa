package errors

import sharedErrors "backend-golang/shared/errors"

// Authentication & Login Errors
var (
	ErrInvalidCredentials   = sharedErrors.Unauthorized("invalid_credentials", "Username atau password salah. Coba lagi!")
	ErrTooManyLoginAttempts = sharedErrors.TooManyRequests("too_many_login_attempts", "Terlalu banyak percobaan login. Coba lagi dalam 15 menit")
	ErrUserInactive         = sharedErrors.Forbidden("user_inactive", "Akun belum aktif. Silakan lakukan verifikasi email")
	ErrAccountLocked        = sharedErrors.Locked("account_locked", "Akun sementara dikunci karena terlalu banyak percobaan login")
)

// Registration & Email Verification Errors
var (
	ErrEmailNotRegistered      = sharedErrors.ValidationError("email_not_registered", "Email belum terdaftar. Silakan lakukan pendaftaran")
	ErrEmailNotFound           = sharedErrors.NotFound("email_not_found", "Email tidak ditemukan")
	ErrEmailAlreadyVerified    = sharedErrors.ValidationError("email_already_verified", "Email sudah diverifikasi")
	ErrEmailNotVerified        = sharedErrors.Forbidden("email_not_verified", "Akun belum aktif. Silakan verifikasi email terlebih dahulu")
	ErrVerificationCodeExpired = sharedErrors.BadRequest("verification_code_expired", "Kode verifikasi sudah kadaluwarsa")
	ErrVerificationCodeInvalid = sharedErrors.BadRequest("verification_code_invalid", "Kode verifikasi tidak valid")
)

// Token Management Errors
var (
	ErrGenerateToken        = sharedErrors.InternalServer("generate_token_failed", "Gagal membuat token akses")
	ErrGenerateRefreshToken = sharedErrors.InternalServer("generate_refresh_token_failed", "Gagal membuat refresh token")
	ErrSaveRefreshToken     = sharedErrors.InternalServer("save_refresh_token_failed", "Gagal menyimpan refresh token")
	ErrTokenExpired         = sharedErrors.Unauthorized("token_expired", "Token sudah kadaluwarsa")
	ErrInvalidRefreshToken  = sharedErrors.Unauthorized("invalid_refresh_token", "Refresh token tidak valid atau sudah dicabut")
)

// Password Reset Errors
var (
	ErrResetTokenExpired = sharedErrors.BadRequest("reset_token_expired", "Token reset password sudah kadaluwarsa")
	ErrResetTokenInvalid = sharedErrors.BadRequest("reset_token_invalid", "Token reset password tidak valid")
	ErrOldPasswordWrong  = sharedErrors.BadRequest("old_password_wrong", "Password lama tidak benar")
)

// User Management Errors
var (
	ErrUserCreationFailed  = sharedErrors.InternalServer("user_creation_failed", "Gagal membuat akun pengguna")
	ErrUserRetrievalFailed = sharedErrors.InternalServer("user_retrieval_failed", "Gagal mengambil data pengguna")
	ErrUserUpdateFailed    = sharedErrors.InternalServer("user_update_failed", "Gagal memperbarui data pengguna")
)

// Session & Security Errors
var (
	ErrSessionExpired   = sharedErrors.Unauthorized("session_expired", "Sesi sudah berakhir. Silakan login ulang")
	ErrInvalidSession   = sharedErrors.Unauthorized("invalid_session", "Sesi tidak valid")
	ErrPermissionDenied = sharedErrors.Forbidden("permission_denied", "Tidak memiliki izin untuk aksi ini")
)
