dockerfillers
=============

The goal of this project is to provide tools that enhance Docker's functionality. The primary utility is "dfimage"
(stands for "docker filler image"), which provides transparency into the content of each image layer (the
files/directories that have been added, updated or removed).

The following commands are supported:

    * diffchanges  (sudo dfimage diffchanges <image name>) - prints changes (file additions,
     modifications and deletions) of the image relative to its ancestors.

    * diffsize (sudo dfimage diffsize <image name>) - prints the overall size contribution
     of the image relative to its ancestors. Note that that might be a negative number, e.g.
     when a big file that exists in the parent image has been deleted in the child image. This
     can give a hint to the flattening potential of the image.

    * help - prints available commands


dfimage currently works by analyzing the data managed by the Docker storage driver, hence it needs to be run with "sudo".

Currently supported storage backends:
    * aufs
    * devicemapper

Build
=====

cd dfimage; go build

Use
===

Usage: dfimage command [arg...]

Commands:
   diffsize          size of an image (in bytes), relative to its ancestors
   diffchanges       changes of an image relative to its ancestors
   help              lists available commands


Don't forget to use sudo to run dfimage (see above for the reason).

Example
========

sudo ./dfimage diffsize busybox:latest

Image id: e72ac664f4f0c6a061ac4ef332557a70d69b0c624b6add35f1c181ff7fff2287 [busybox:latest]:
0

Image id: e433a6c5b276a31aa38bf6eaba9cd1cfd69ea33f706ed72b3f20bafde5cd8644 []:
A /bin 4096
A /bin/ash 7
A /bin/busybox 561992
A /bin/cat 7
A /bin/catv 7
A /bin/chattr 7
A /bin/chgrp 7
...
2593047

Image id: df7546f9f060a2268024c8a230d8639878585defcc1bc6f79d2728a13957871b []:
0
