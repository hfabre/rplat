package game

func RemoveIndexStar(s []Star, index int) []Star {
	return append(s[:index], s[index+1:]...)
}
