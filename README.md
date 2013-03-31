goblog-prototype
================

Course project to build a blogging engine in Go.

# Usage

To start the blog in dev mode, you need the following variables set:
```
export PORT=5000
export DATABASE_URL="<a url to your postgres url>"
export DEBUG=true
```

The `DEBUG` variable will generate random data in the database.  If you
do not want that, simply let the variable an empty string.
