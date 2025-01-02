# File Servers and CDNs

## File Storage

Building a (good) web application almost always involves handling "large" files of some kind - whether its static images and videos for a marketing site, or user generated content like profile pictures and video uploads, it always seems to come up.

### Large Files

You're probably already familiar with small structured data; the stuff that's usually stored in a relational database like Postgres or MySQL. 
I'm talking about simple, primitive data types like:

* `user_id` (integer)
* `is_active` (boolean)
* `email` (string)

Large files, or "large assets", on the other hand, are giant blobs of data encoded in a specific file format and measured in kilo, mega, or gigabytes. 

### Database

SQLite is a traditional relational database that works out of a single flat file, meaning it doesn't need a separate server process to run.

```sh
brew install sqlite3
```

### Videos

"Videos" have 3 things to worry about:

* Metadata: The title, description, and other information about the video
* Thumbnail: An image that represents the video
* Video: The actual video file

### Multipart Uploads

So you might already be familiar with simple JSON/HTML form POST requests. 
That works great for small structured data (small strings, integers, etc.), but what about large files?

We don't typically send massive files as single JSON payloads or forms. 
Instead, we use a different encoding format called multipart/form-data. 
In a nutshell, it's a way to send multiple pieces of data in a single request and is commonly used for file uploads. 
It's the "default" way to send files to a server from an HTML form.

### Encoding

But we can't store an image in a SQLite column?... Right?

To do so, we can actually encode the image as a base64 string and shove the whole thing into a text column in SQLite. 
Base64 is just a way to encode binary (raw) data as text. 
It's not the most efficient way to do it, but it will work for now.

### Using the Filesystem

It's usually a bad idea to store large binary blobs in a database, there are exceptions, but they are rare. 
So what's the solution? 
Store the files on the file system. 
File systems are optimized for storing and serving files, and they do it well.

### Mime types

There are an infinite number of things we could consider "large files". 
But within the context of web development, the most common types of large files are probably:

1. Images: PNGs, JPEGs, GIFs, SVGs, etc.
2. Videos: MP4s, MOVs, AVIs, etc.
3. Audio: MP3s, WAVs, etc.
4. Static web templates: HTML, CSS, JS, etc.
5. Administrative files: PDFs, Word docs, etc.

A mime type is just a web-friendly way to describe format of a file. 
It's kind of like a file extension, but more standardized and built for the web.

Mime types have a type and a subtype, separated by a /. 
For example:

* image/png
* video/mp4
* audio/mp3
* text/html

When a browser uploads a file via a multipart form, it sends the file's mime type in the Content-Type header.

### Live Edits

If a user were able to live edit a file (think Google Docs or Canva) we'd have to approach our storage problem differently. 
We wouldn't just be managing new versions of "static" files, we would need to handle every tiny edit (keystroke) and sync updated changes to our server. 

## Caching

When a user visits a web application for the first time, their browser downloads all the files required to display the page: HTML, CSS, JS, images, videos, etc. 
It then "caches" (stores) them on the user's machine so that next time they come back, it doesn't need to re-download everything. 
It can use the locally stored copies.

### Cache Busting

Browsers cache stuff for good reason: it makes the user experience snappier and, if the user is paying for data, cheaper.

That said, sometimes (like in the last lesson) we don't want the browser to cache a file - we want to be sure we have the latest version. 
One trick to ensure that we get the latest is by "busting the cache". 
A simple tactic is to change the URL of the file a bit. 

To cache bust, we want to alter the URL so that:

* The browser thinks its a different file
* The server thinks its the same file

### Cache Headers

Query strings are a great way to brute force cache controls as the client - but the best way (assuming you have control of the server, and c'mon, we're backend devs), is to use the Cache-Control header. 
Some common values are:

* `no-store`: Don't cache this at all
* `max-age`=3600: Cache this for 1 hour (3600 seconds)
* `stale-while-revalidate`: Serve stale content while revalidating the cache
* `no-cache`: Does not mean "don't cache this". It means "cache this, but revalidate it before serving it again

When the server sends Cache-Control headers, its up to the browser to respect them, but most modern browsers do.

### New Files

"Stale" files are a common problem in web development. 
And when your app is small, the performance benefits of aggressively caching files might not be worth the complexity and potential bugs that can crop up from not handling cache behavior correctly.

## AWS S3

### Single Machine

In a "simple" web application architecture, your server is likely a single machine running in the cloud. 
That single machine probably runs:

1. An HTTP server that handles the incoming requests
2. A database running in the background that the HTTP server talks to
3. A file system the server uses to directly read and write larger files

### AWS

AWS (Amazon Web Services) is one of the (at least in my mind) "Big Three" cloud providers. 
The other two are Google Cloud and Microsoft Azure.

AWS is the oldest, largest, and most popular of the three, generally speaking.

### Serveless

"Serverless" is an architecture (and let's be honest, a buzzword) that refers to a system where you don't have to manage the servers on your own.

You'll often see "Serverless" used to describe services like AWS Lambda, Google Cloud Functions, and Azure Functions. 
And that's true, but it refers to "serverless" in its most "pure" form: serverless compute.

AWS S3 was actually one of the first "serverless" services, and is arguably still the most popular. 
It's not serverless compute, it's serverless storage. 
You don't have to manage/scale/secure the servers that store your files, AWS does that for you.

Instead of going to a local file system, your server makes network requests to the S3 API to read and write files.

### Upload Object

### Architecture

### SDK and S3

## Object Storage

## Video Streaming

## Security

## CDNs

## Resiliency