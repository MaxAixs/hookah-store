package messages

const (
	UserRegistered = "user registration successfully"
	PasswordReset  = "password reset successfully"
)

const (
	singUpEvenType         = "sign_up"
	resetPasswordEventType = "reset_password"
)

var MapMsg = map[string]string{
	singUpEvenType:         UserRegistered,
	resetPasswordEventType: PasswordReset,
}
