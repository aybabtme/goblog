goblog-prototype
================

Originally a course project where we had to code anything that uses a SQL database, we decided to take the opportunity to learn Go and write the whole thing from scratch.

# Notes

__As it stands, the URL of OAuth callbacks is hardcoded to `flying-unicorn.aybabt.me:5000`__.
 I might change that sometime in the future to use a config file.  Same thing for the Google API secrets, they're hardcoded in there.
 
To run the blog locally, edit your `/etc/hosts/` file and add the following line:

```
127.0.0.1   flying-unicorn.aybabt.me
```

# Usage

To start the blog in dev mode, you need the following variables set:

```
export PORT=5000
export DATABASE_URL="<a url to your postgres url>"
```

Then start the blog:

```
go get github.com/aybabtme/goblog
goblog
```

This will start the blog with no users and no data in it.  You might want to have an author, so that you can write posts.

To do this, use the `--create-admin` flag.  When you start the blog, it will askes you to do somethings to create
an author user, and then the blog will start per-se.

```
goblog --create-admin
```

If you'd like to get some random lorem ipsum data in your blog, run it using the `--debug` flag.

```
goblog --debug --create-admin
```

# Writing posts and comments

Comments and posts are converted to HTML using a Markdown compiler.  The syntax is kind-of Github-like.  Any HTML you leave in there will be escaped.
