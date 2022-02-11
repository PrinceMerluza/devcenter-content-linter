package main

func main() {
	paramsData, err := getParams()

	if err != nil {
		exitGracefully(err)
	}

	cloneBlueprint(paramsData.repoPath)
}
