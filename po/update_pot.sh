POT_FILE=deezer-scope.pot
cp ../deezer/deezer-scope.labsin_deezer-scope.ini.in ./
cp ../deezer/deezer-scope.labsin_deezer-scope-settings.ini.in ./
xgettext --keyword=Gettext --package-name=deezer-scope --package-version=0.1 --msgid-bugs-address=sam.sgrs@gmail.com -Lc -o$POT_FILE --from-code="utf-8" -cTRANSLATORS ../deezer/*.go
intltool-extract --type=gettext/ini --update deezer-scope.labsin_deezer-scope.ini.in
intltool-extract --type=gettext/ini --update deezer-scope.labsin_deezer-scope-settings.ini.in
intltool-extract --type=gettext/xml --update ../myapps.xml.in
xgettext -c --keyword=N_ --package-name=deezer-scope --package-version=0.1 --msgid-bugs-address=sam.sgrs@gmail.com -Lc -j -o$POT_FILE --from-code="utf-8" *.ini.in.h
xgettext -c --keyword=N_ --package-name=deezer-scope --package-version=0.1 --msgid-bugs-address=sam.sgrs@gmail.com -Lc -j -o$POT_FILE --from-code="utf-8" ../myapps.xml.in.h
for PO_FILE in *.po; do
    msgmerge -U $PO_FILE $POT_FILE
done
rm deezer-scope.labsin_deezer-scope.ini.in
rm deezer-scope.labsin_deezer-scope-settings.ini.in
rm *.ini.in.h
rm ../myapps.xml.in.h
