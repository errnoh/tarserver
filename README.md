tarserver
=========

Because someone mentioned web servers that would serve gzipped tar files..

Usage:
------

    curl -F file=@file1.go -F file=@file2.txt http://127.0.0.1:8080/tar/gz > tarredstuff.tar.gz

..or just send files as MIME multipart ( RFC 2046 ) with some other software..
..or don't send files, just send some random data (though skip /tar/ if you do that)..

Supported functions:
--------------------

/tar - tar files that are sent as MIME multipart form.
/gzip or /gz/ - gzip data
/zlib - zlib data
/echo - echo server

Paths can be combined:
----------------------

/tar/gz - first tar the files and then gzip them

TODO:
-----

* https support
