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

func BuildObjectPathWithKey(userSub, album, key string) string {
	albumPath := BuildObjectPath(userSub, album)
	return *albumPath + key
}
