#!/bin/sh -e

TMPFILE=tmp.cov
OUTFILE=profile.cov

echo "mode: count" > $OUTFILE

for pkg in $(go list ./...); do
    go test -covermode count -coverprofile $TMPFILE $pkg
    if [ -f $TMPFILE ]; then
        cat $TMPFILE | sed 1d >> $OUTFILE
        rm $TMPFILE
    fi
done
