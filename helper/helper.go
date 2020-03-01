package helper

func AddSlash(str string) *string {
	temp := str + "/"
	return &temp
}

func BuildObjectPath(userSub, album string) *string {
	temp := userSub + "/"
	if album != "" {
		temp = temp + album + "/"
	}
	return &temp
}
