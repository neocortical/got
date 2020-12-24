# got
A clone of the git vcs written in Go. 

This work is based on the book [_Building Git_](https://shop.jcoglan.com/building-git/) by James Coglan. This book is excellent and I highly recommend it!

I will try to tag the repository state according to the chapter titles.

## About The Repo
This year I've been working on improving my understanding of git. Like many developers, I've long survived on just enough git proficiency to scraped by. But I wanted to get better. 

The first thing I did was spend some time learning interactive rebasing. That was a big leap forward. I also bought Katie Sylor-Miller and Julia Evans' zine [_Oh Shit, Git!_](https://wizardzines.com/zines/oh-shit-git/). It's a great primer and cookbook for working with git, and fun to read as well!

But I wanted to really understand git's internals better, and I also like coding projects, so when I saw James's book, I bought it immediately. 

I'm not a huge fan of Ruby, and I work professionally in Go, so I thought it would be really fun to read the book and implement git in Go as I read. An advantage of doing this is that I don't have to resist the urge to just copy/paste the code examples into my project — they ain't gonna compile. And Go and Ruby are very different languages, so there really wasn't a temptation to just transliterate James's Ruby code into Go.

This process has been really challenging and extremely rewarding. I'm less than a third of the way through the book (unless this README has gone stale) and I already feel like I've leveled up my understanding of git internals multiple times. In addition, I'm learning aspects of Go that I haven't encountered personally or professionally, so win win! 

I am also interested in (finally) learning Rust, so I'm eyeing a future project of doing this all again in a new language! :D

## Install
`go get -u github.com/neocortical/got`

## Use
Just like using `git`, except it's buggy and most of the features aren't implemented!