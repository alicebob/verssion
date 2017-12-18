Atom feeds of (stable) software versions, as found in Wikipedia.

Setup
=====

    postgres@yourmachine:~$ createdb -O youruser verssion
    youruser@yourmachine:~/verssion/$ make db
    youruser@yourmachine:~/verssion/$ make && ./cmd/web/web -base https://yourwebsite.example

`make integration` will use the `verssion` database, and wipe everything from
it. Just so you know.

&c.
===

[![Build Status](https://travis-ci.org/alicebob/verssion.svg?branch=travis)](https://travis-ci.org/alicebob/verssion)

`cmd/web/web` is the main HTTP server. `cmd/wikimon/wikimon` is a helper to
compare verssion and wikipedia against the websites of projects. I run it via cron
to help keep wikipedia up to date.
