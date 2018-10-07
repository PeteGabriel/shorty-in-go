package main

func main() {
	a := App{}
	//TODO this should go into env vars
	a.Initialize("dummy0", "dummy1", "dummy2")
	a.Run(":8080")

}
