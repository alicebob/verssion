Atom feeds of (stable) software versions, as found in Wikipedia.

Bugs
====

* Versions on the wikipedia 'Firefox' page are not supported.
* 'Foo_bar' page is the same page as 'Foobar', but wikipedia doesn't do a proper redirect.
* Better name pun.

Setup
=====

Postgres
--------

postgres@yourmachine:~$ createdb -O youruser w
youruser@yourmachine:~/w/$ make db
