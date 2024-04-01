package helper

import ()

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
