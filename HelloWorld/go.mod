module HelloWorld

go 1.20

require mypackage v0.0.0

require github.com/q1mi/hello v0.1.1 // indirect

replace mypackage => ../mypackage
