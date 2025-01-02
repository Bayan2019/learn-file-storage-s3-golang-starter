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

### New Files

## AWS S3

## Object Storage

## Video Streaming

## Security

## CDNs

## Resiliency