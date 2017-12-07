Atom feeds of (stable) software versions, as found in Wikipedia.

Bugs
====

* 'Foo_bar' page is the same page as 'Foobar', but wikipedia doesn't do a proper redirect.
* Better name pun.

Setup
=====

    postgres@yourmachine:~$ createdb -O youruser verssion
    youruser@yourmachine:~/verssion/$ make db
    youruser@yourmachine:~/verssion/$ make && ./cmd/web/web -base https://yourwebsite.example

`make integration` will use the `verssion` database, and wipe everything from
it. Just so you know.
