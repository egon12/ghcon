package cover

var exampleUse = false

func NotCoverFuncExamples() {
	str := "Hello"
	str += "\nMy name is"
	str += "\nSecret"

	if exampleUse {
		str += "\n But this is"
		str += "\n Something that should not be test"
		println(str)
	}
}
