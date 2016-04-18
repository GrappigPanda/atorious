[![Build Status](https://travis-ci.org/GrappigPanda/notorious.svg?branch=devel)](https://travis-ci.org/GrappigPanda/notorious) [![Go Report Card](https://goreportcard.com/badge/github.com/GrappigPanda/notorious)](https://goreportcard.com/report/github.com/GrappigPanda/notorious)
# Notorious
Hello everybody! Notorious aims to be a highy extensible tracker implemented in golang. Right now, Notorious uses Redis to store peer information for quick retrieval and to alleviate some of the burden of storing peer information from the tracker itself. Moreover, some of the core functionalities which Notorious hopes to gain include the following:
```
1. Ratio tracking using a SQL backend (I will be using an ORM layer so that most SQL DBs will be supported).
2. "Complete" Ratioless tracking: no peer information is ever stored
3. "Semi-ratioless" tracking: we store just statistics information + "Complete" tracking info.
4. "Normal" tracking: a unique key is assigned to a peer and we store seeds/leeches, downloaded bytes/uploaded bytes.
5. Automatic Redis docker deployment. I like docker and Redis is my peer-storage of choice and they work well together.
6. Speed and scalability are always in the back of my mind, even if my decisions don't always reflect that
```


There's probably a lot more! Check out my [issues page](https://github.com/GrappigPanda/notorious/issues)

# Deployment

Deployment is one of the funnest things from this project for me so far because I've gotten to use a lot of cloud technologies which I typically wouldn't get to deal with. As of the moment, I'm using [Docker Cloud](https://cloud.docker.com/) and deploying to a 5 $ [Digital Ocean](https://m.do.co/c/39961c9b71bf) droplet. As you'll see, there's a dockerfile in my repo as well as a supervisord.conf and to build it, you'd want to run:
```
go build main.go
docker build -t notorious .
```
and then you can run
```
docker run <-d> -p 3000:3000 notorious (where the <> indicate the -d is optional, it starts it in daemon mode)
```
Et voilà, you have a copy of notorious running,.

This will build the docker image which you can then run it either on local bare-metal, or if you're interested in the Docker Cloud route (which I highly recommend) head over to their [Documentation](https://docs.docker.com/docker-cloud/getting-started/) They'll explain deployment 1 000% better than I can.

## Please note:
I have not yet included instructions for running everything on bare-metal and I'm using docker because it just works (sorry, I try to keep the buzzwordiness down [despite deploying notorious with mariadb, redis, docker, go...]), but those are soon to come. I just honestly thing the docker way is 10x (just like me, heh) easier to deploy with at this particular moment.

# Test it out?

Hey, while you're here why don't you try out NetBSD after testing out my tracker [NetBSD-7.0.arm64.torrent](https://github.com/GrappigPanda/notorious/blob/master/NetBSD-7.0-amd64.iso.torrent)
(I'm not part of the NetBSD project or anything, I just like the product -- also, you probably run Intel, so this torrent doesn't really do you any good unless you're looking to virtualize NetBSD version 7.0. So for that one person out there, this one's for you).

# Contact Notorious

Notorious is a project which I've had a ton of fun learning Go in, but do realize I'm still learning Go so I do make non-idiomatic decisions. If you see anywhere that you think I could improve my code or golang usage, please:

[open an issue](https://github.com/GrappigPanda/notorious/issues/new)

[tweet me](http://twitter.com/GrappigPanda)

[or email me](mailto:ian@ianleeclark.com)

# License
The MIT License (MIT)

Copyright (c) 2016 Ian Clark

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
