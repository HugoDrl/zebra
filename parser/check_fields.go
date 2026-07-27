package parser


func checkLogValidity(log *Log, settings *ParseSettings) bool {
	if !settings.StartDate.IsZero() && log.Time.Compare(settings.StartDate) < 0 {
		return false
	}
	if !settings.EndDate.IsZero() && log.Time.Compare(settings.EndDate) > 0 {
		return false
	}
	if settings.Level != "" && log.Level != settings.Level {
		return false
	}
	if settings.Service != "" && log.Service != settings.Service {
		return false
	}
	return true
}
