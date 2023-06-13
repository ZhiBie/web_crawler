
From: https://realpython.com/async-io-python/#a-full-program-asynchronous-requests

In this section, you’ll build a web-scraping URL collector, areq.py, using aiohttp, a blazingly fast async HTTP client/server framework. (We just need the client part.) 

The high-level program structure will look like this:

1. Read a sequence of URLs from a local file, urls.txt.

1. Send `GET` requests for the URLs and decode the resulting content. If this fails, stop there for a URL.

1. Search for the URLs within `href` tags in the HTML of the responses.

1. Write the results to `foundurls.txt`.

1. Do all of the above as asynchronously and concurrently as possible. (Use `aiohttp` for the requests, and `aiofiles` for the file-appends. These are two primary examples of IO that are well-suited for the async IO model.)

The second URL in the list should return a 404 response, which you’ll need to handle gracefully. If you’re running an expanded version of this program, you’ll probably need to deal with much hairier problems than this, such a server disconnections and endless redirects.

The requests themselves should be made using a single session, to take advantage of reusage of the session’s internal connection pool.
