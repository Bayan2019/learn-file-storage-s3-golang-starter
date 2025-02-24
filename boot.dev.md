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

A few things to think about:

1. Who can access your bucket, and which parts of your bucket can they access?
2. What actions can they take?
3. How are they authenticated? 
    And from where can they authenticate?

While it's great that an attacker would need to steal your AWS credentials to be able to maliciously change the contents of your bucket, relying only on the secrecy of keys is often not enough.

Keys and passwords are compromised all the time.

One way to add an additional layer of security is to ensure that your keys can only be used from certain (virtual) locations. 
Then an attacker would need your keys and to be on your network to gain access.

### Application vs. Developer

The right thing to do in production is to have your code running on a dedicated server - not on your home desktop.

### Scoping Permission

A critical rule of thumb in cyber security is the principle of least priviledge: You should allow the fewest permissions possible that can still get the job done.

For example, your user is in the "manager" group which we gave "full admin access" to. 
Especially at smaller companies, it's common for folks to have more permissions than they truly need, usually for the sake of speed and convenience.

But that's not the most secure way to do things.

### Roles

Now that we have a policy for the tubely app, we need a role to attach it to.

### Private Bucket

Public buckets are useful when you want to serve public content directly from them, like user profile pictures, for example. 
However, you should only use them when you're certain all the content should be public, and you're okay with the risks of anyone on the internet using the bandwidth you pay AWS for to download your assets over and over again...

A good use case for a public bucket might be:

* Users' profile pictures
* Public certificates of completion (we do this for Boot.dev!)
* Dynamically generated images for social sharing (like the link previews you see on Twitter)

While a **private bucket** might contain:

* A user's privately uploaded documents
* A user's draft content that they haven't published yet
* The org's video content that's only available to paying customers

### Signed URLs

Presigned URLs are a way to give temporary access to a private object in S3. 
S3 will generate a URL (by attaching a cryptographic signature) that allows access to the object for a limited time. 
To be clear, it doesn't require the user to be logged in - it's just a URL that expires.

The idea is that we'll generate these URLs with very short life spans, and then only give them to users who have already been authenticated by your application.

### Encryption

Although our S3 bucket is private (which means outsiders can't gain access to the files directly without credentials), it's still good for stuff to be encrypted. 
After all, what if a hacker physically walked into the data center to read our customers' secrets directly?

#### At Rest

Files in S3 are encrypted at rest ("at rest" just means "while they're sitting in storage on disk") by default. 
This was not always the case, but it is now! 
You don't need to do anything, the S3 service takes care of all of that for you. 
When you access S3 with your credentials, the service decrypts the files for you before handing them over.

#### In Transit

When you're uploading or downloading files from S3, how do you know that someone can't intercept the data as it travels through the internet? 
Well, when you access S3 via the web, you're using httpS. 
The S means that the data is encrypted as it travels between your computer and the S3 service.

When you access S3 via the SDK (in your Go code), it also uses HTTPS by default. 
So as long as you don't go out of your way to disable encryption, you're good to go.

## CDNs

### Regions

A region is a geographic location where AWS has data centers. 
Data centers are clustered into "availability zones" (or "AZ" for acronym masochists enjoyers).

By default, your S3 bucket is replicated across multiple availability zones in a single region. 
That said, there are options to automatically replicate your bucket's data across multiple regions. 
(We won't do that, but it's good to know about.)

### CDNs

A Content Delivery Network (CDN) is a (typically global) network of servers that caches and delivers content to users based on their geographic location.

When we give users a URL to an S3 object, they'll download that object from the S3 service in the region that our bucket lives in (for me, that's us-east-2, near Ohio in the USA).

If a user in Australia tries to download that object, they're going to have to wait for the data to travel from Ohio to Australia... and that's a long way! 
A CDN, like AWS CloudFront, can help with that. 
It takes a static asset, like an image or video, and caches it on servers all over the world. 
When a user requests the asset, they get it from the server closest to them, which is much faster.

### Use CloudFront

Signed URLs are useful for truly private content, but if all you need is more protection and control over files that you want to make publicly accessible, a CDN is a better choice. 
CDN's like CloudFront not only offer better security than serving files directly from S3 (due to more granular controls, firewalls, and DDoS protection), but they also offer better performance.

### Invalidations

A CDN is a massive, globally distributed cache. 
Sure, we get massive performance improvements, because users that are geographically close to an edge server can download assets much faster than if they had to travel to the origin server.

But what happens when we update an asset? 
How long does it take the edge servers to update their versions? 
The answer is: it depends. 
That's always the tradeoff with cache - you need to deal with invalidations. 
Luckily CloudFront makes it fairly easy to force invalidations of the cache.

An invalidation is a request to remove an object from the cache. 
That means the next time a user requests the object, the edge server will have to go back to the origin server to get the latest version. 
That means it will be slower for the first user, but fast again for subsequent users.

### Why CDNs?

A CDN like CloudFront has two purposes (as far as the context of this course is concerned):

* **Speed**: Users get content from the server closest to them, which is faster than getting it from the origin server.
* **Security**: The origin server is hidden from the public internet, and only the CDN can access it. 
    This is a security measure that can help prevent DDoS attacks and other malicious activity.

Some CDNs, like CloudFlare, (not to be confused with CloudFront) are known for their incredibly robust security features. 
Things like DDoS protection, Web Application Firewalls, etc.

Images and videos are certainly common, but in reality any static asset is a good fit for a CDN. 
Here at Boot.dev, we use CloudFlare's CDN to serve the static assets for our frontend:

* Images
* HTML
* CSS
* JS

We deploy on their edge network, which means that our users get the initial HTML document quickly. 
That said, our backend server is a Go application running in a single region in the United States, so any dynamic requests to our API still have to come all the way back to the US.

## Resiliency

### Availability

I mentioned earlier that one of the big advantages of "serverless" (and in particular, S3) is that it takes care of a lot of the "IT ops" work that traditionally engineers at every company had to homebrew.

One of those is availability: how often your service is up and running, serving user requests. 
It's often measured in "nines" - like "three nines" (99.9%) or "five nines" (99.999%).

See, users don't like when they log into your web app and stuff isn't loading. 
They don't like to hear that you're "down for maintenance".

AWS and S3 aren't perfect - but they are really good at availability. 
When AWS has outages, it's big news. 
Partly because so much of the internet runs on AWS, but also because they're rare.

You could build your own cluster of servers with better than or equal to one of the large cloud provider's availability. 
But its very hard, and very expensive.

### Reliability

Okay, so let's say you've got 5 9's of availability. 
That's 99.999% uptime. 
Pretty solid. 
But what about reliability?

Reliability is about how well your system works when it's up. 
For example, maybe your server is responding to HTTP requests, but it's returning erroneous data because some dependency is down. 
That's not reliable.

The reliability of S3 is very high out of the box.

### Durability

Durability is the last in the resiliency trifecta:

* Availability
* Reliability
* Durability

It's about how well your data survives in the event of an outage. 

For example, let's say you're running your own single server:

1. What happens if the intern accidentally rm -rfs the user_pics directory?
2. What happens if the server's hard drive fails?
3. What happens if the data center it's in catches fire?

These are all durability questions. 
Durability is primarily about backups and redundancy. 
In the case of S3, it automatically replicates your data across multiple servers. 
If one goes down, the backups are there.

According to these docs S3's standard storage provides 99.999999999% durability and 99.99% availability of objects over a given year. Nice.

### Bucket Versioning

By default, S3 does not store multiple versions of an object. 
If you upload a file to a key that already contains an object, the old object is overwritten.

Bucket versioning is an optional feature where the bucket stores multiple versions of an object. 
It helps:

* Prevent accidental deletion
* Rollback to previous versions of files
* Store multiple versions of files in the same key
