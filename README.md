goblog
================

Originally a course project where we had to code anything that uses a SQL database, we decided to take the opportunity to learn Go and write the whole thing from scratch.

# Notes

__As it stands, the URL of OAuth callbacks is hardcoded to `flying-unicorn.aybabt.me:5000`__.
 I might change that sometime in the future to use a config file.  Same thing for the Google API secrets, they're hardcoded in there.  They are deactivated, so you need to replace them with your own.
 
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

# Known bugs

* _Template rendering during concurrent connections._ The way templates are rendered by the Controllers is not thread safe.  When two or more goroutine meet the same template variable during execution, they may conflict with one another and result in a broken pipe, which resets the connection.  A fix for this would be to offer the Controllers a `chan *template.T` instead of just a `*template.T`.  The chan would contain `runtime.NumCPU()` templates and every controller calling a template would remove one from the chan, render with the template they took then put the template back into the channel.  Since `GOMAXPROCS` is set to `NumCPU()`, this would not result in any slowdown.  Doing so could also allow for live changes to the templates, having a watching goroutine that looks up for changes in the template files and replace the templates in the chan by new versions.
* _Erroneous credential handling._ There seems to be a bug when two users log into the website from the same computer but using different Google accounts.  The firstly used account seems to be the only one connecting, such that when you try to connect using another Google account, your get logged into the blog as the prior user.  The issue has not been further investigated.
