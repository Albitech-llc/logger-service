package interfaces

type Service interface {
	LogInfo(service string, message string) error
	LogError(service string, message string) error
	LogWarning(service string, message string) error
	LogMessage(service string, message string, level string) error
}
