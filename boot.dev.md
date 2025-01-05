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

S3 is really simple.

File A goes in bucket B at key C. That's it. 
You only need 2 things to access an object in S3:

* The bucket name
* The object key

Buckets have globally unique names - If I make a bucket called "bd-vids", you can't make a bucket called "bd-vids", even if you're in a separate AWS account. This makes it really easy to think about where your data lives.

### SDKs and S3

An SDK or "Software Development Kit" is just a collection of tools (often involving an importable library) that helps you interact with a specific service or technology.

AWS has official SDKs for most popular programming languages. 
They're usually the best way to interact with AWS services. 

When you as a human interact with AWS resources, you'll typically use the web console (GUI) or the CLI. 
When your code interacts with AWS resources, you'll use the SDK within your code.

## Object Storage

If you squint really hard, it feels like S3 is a file system in the cloud... but it's not. 
It's technically an object storage system - which is not quite the same thing.

Object storage is designed to be more scalable, available, and durable than file storage because it can be easily distributed across many machines:

* Objects are stored in a flat namespace (no directories)
* An object's metadata is stored with the object itself

### File System Illusion

Directories are really great for organizing stuff. 
Storing everything in one giant bucket makes a big hard-to-manage mess. 
So, S3 makes your objects feel like they're in directories, even though they're not.

Keys inside of a bucket are just strings. And strings can have slashes, right? Right.

### Dynamic Path

Schema architecture matters in a SQL database, and prefix architecture matters in S3. 
We always want to group objects in a way that makes sense for our case, because often we'll want to operate on a group of objects at once.

If you don't have any prefixes (directories) to group objects, you might find yourself iterating over every object in the bucket to find the ones you care about. 
That's slow and expensive.

## Video Streaming

### Streaming

Now that (almost) no one is on dial-up 256k modems, we typically don't worry about "streaming" smaller files like images.

But giant audio files (like audio books), and especially large video files should be streamed rather than downloaded. 
At least if you want your user to be able to start consuming the content immediately.

The simplest way to stream a video file on the web (imo) is to take advantage of two things:

1. The native HTML5 `<video>` element. 
    It stream video files by default as long as the server supports it.
2. The Range HTTP header. 
    It allows the client to request specific byte ranges of a file, enabling partial downloads. 
    S3 servers support it by default.

### MP4

### More Requests

### Other approaches

Because we're mostly concerned with S3 and file storage in this course, we won't be doing a deep dive on all the file formats for video streaming. That said, I want to briefly cover a few of the most common ones, and point out when you might need to ditch a "plain" .mp4 file.

1. Adaptive streaming: Standard mp4 files have a single resolution and bitrate. 
    If a user's connection speed is unstable, HLS or MPEG-DASH allows for changing the quality of the stream on the fly. 
    You may have noticed on YouTube or Netflix that your video quality changes based on your connection speed. 
    Dropping to lower resolution is better than endlessly buffering.
2. Live streaming: Standard mp4 files are not designed to be updated in real-time. 
    You'd want to use a lower-latency protocol like WebRTC or RTMP for live streaming.

## Security

## CDNs

## Resiliency