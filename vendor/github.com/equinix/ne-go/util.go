package ne

//String returns pointer to a given string value
func String(s string) *string {
	return &s
}

//StringValue returns the value of a given string pointer
//or empty string if the pointer is nil
func StringValue(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

//Int returns pointer to a given int value
func Int(i int) *int {
	return &i
}

//IntValue returns the value of a given int pointer
//or 0 if the pointer is nil
func IntValue(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}

//Bool returns pointer to a given bool value
func Bool(b bool) *bool {
	return &b
}

//BoolValue returns the value of a given bool pointer
//or false if the pointer is nil
func BoolValue(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}
