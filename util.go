package aferocopy

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func ignore(error) {}
