Monitor (stable) software versions via Wikipedia

Bugs
====

* 'Firefox' versions on wikipedia are not supported.
* 'Foo_bar' page is the same as 'Foobar', but it's not a proper redirect.

Setup
=====

Postgres
--------

postgres@yourmachine:~$ createdb -O youruser w
youruser@yourmachine:~/w/$ make db
