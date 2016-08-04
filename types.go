package main

type Configuration struct {
	Login struct {
		Username string `json:username`
		Password string `json:password`
	}
	Paths struct {
		PathToHD string `json:pathToHD`
	}
}
