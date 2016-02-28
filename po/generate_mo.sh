#!/bin/bash
for PO_FILE in *.po; do
    LANG=${PO_FILE:0:(-3)}
    mkdir -p "../share/locale/$LANG/LC_MESSAGES/"
    MO_FILE="../share/locale/$LANG/LC_MESSAGES/deezer-scope.mo"
    msgfmt -c -v -o $MO_FILE $PO_FILE
done
