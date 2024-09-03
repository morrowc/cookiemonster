package main

import "text/template"

const (
	index = `<html><head><title>Index Page</title></head>
<body>
Three options:<br>
<ol>
<li> <a href="/login">Login</a>
<li> <a href="/logout">Logout</a>
<li> <a href="/stuff">Stuff</a>
</ol>
</body>
</html>`

	stuff = `<html><head><title>Stuff Here!</title></head>
<body>Stuff Is here!</body></html>`

	login = `<html><head><title>Login on this page</title></head>
<body>
<form id="loginPage" method="post" action="/login">
<label for="username">UserName: </label><input type="text" id="username"><br>
<label for="passwd">Password: </label><input type="text" id="passwd"><br>
<input type="submit" value="Submit">
</form>
</body>
</html>`

	logout = `<html><head><title>Logout!</title></head>
<body>Logged out!</body></html>`
)

var (
	indexTmpl  = template.Must(template.New("index").Parse(index))
	stuffTmpl  = template.Must(template.New("stuff").Parse(stuff))
	loginTmpl  = template.Must(template.New("login").Parse(login))
	logoutTmpl = template.Must(template.New("logout").Parse(logout))
)
