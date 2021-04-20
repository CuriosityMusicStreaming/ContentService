package app

func assertErr(err error) {
	if err != nil {
		panic(err)
	}
}
