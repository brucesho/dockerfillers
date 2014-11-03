dockerfillers
=============

The goal of this project is to provide tools that enhance Docker's functionality. The primary utility is "dfimage"
(stands for "docker filler image"), which supports multiple commands:

    * diffchanges  (sudo dfimage diffchanges <image name>) - prints changes (file additions,
     modifications and deletions) of the image relative to its ancestors.

    * diffsize (sudo dfimage diffsize <image name>) - prints the overall size contribution
     of the image relative to its ancestors. Note that that might be a negative number, e.g.
     when a big file that exists in the parent image has been deleted in the child image. This
     can give a hint to the flattening potential of the image.


dfimage currently works by analyzing the data managed by the Docker storage driver, hence it needs to be run with "sudo".

Currently supported storage backends:
    * aufs

Next to be supported storage backends:
    * devicemapper
