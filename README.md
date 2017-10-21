Monitor software versions on Wikipedia

Postgres
========
postgres@yourmachine:~$ createdb -O youruser w
youruser@yourmachine:~/w/$ make db

URLs
====

https://www.mediawiki.org/wiki/API:Main_page

https://en.wikipedia.org/w/api.php
format=json
titles=pagea|pageb
rvprop=content
action=query

parsetree:
https://en.wikipedia.org/w/api.php?format=json&page=PostgreSQL&action=parse&prop=parsetree

