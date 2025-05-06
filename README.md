# URL Shortener Service

### Description:

Develop a URL shortener that generates short links for long URLs,
redirects users, and tracks click counts. Store data in a simple database like Redis or a file.


### Skills Showcased:

HTTP handling, database interaction, concurrency, basic web app logic.

### Features:

API endponts:

POST /shorten ( body form data with a key "url" and value of it )
GET /:shortCode ( retrieve original url, and track click counts )

Track and display click counts.

### Tools/Libraries:

- net/http
- local storage

### TODO:

- [x] mux server with 2 endpoints ( POST /shorten, GET /:shortenCode)
- [x] provide a simple openapi document
