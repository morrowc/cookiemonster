# cookiemonster
Simple golang http service that handles authentication using context and cookies.

## Build a simple http server

Server should listen on a designated port (via arg).
Server should accept authentication for (initially) a static userid/passwd.
Server should provide 3 endpoints/urls:
  * /login - accept login details, and fill auth cookie
  * /stuff - send back basic content which is inaccessible without cookie.
  * /logout - destroy the auth cookie.

##
