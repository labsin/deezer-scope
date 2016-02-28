cp ../deezer/deezer-scope.labsin_deezer-scope.ini.in ./
cp ../deezer/deezer-scope.labsin_deezer-scope-settings.ini.in ./
cp ../myapps.xml.in ./
LC_ALL=C intltool-merge -d -u ./ deezer-scope.labsin_deezer-scope.ini.in deezer-scope.labsin_deezer-scope.ini
LC_ALL=C intltool-merge -d -u ./ deezer-scope.labsin_deezer-scope-settings.ini.in deezer-scope.labsin_deezer-scope-settings.ini
LC_ALL=C intltool-merge -x -u ./ myapps.xml.in myapps.xml
mv deezer-scope.labsin_deezer-scope.ini ../deezer/deezer-scope.labsin_deezer-scope.ini
mv deezer-scope.labsin_deezer-scope-settings.ini ../deezer/deezer-scope.labsin_deezer-scope-settings.ini
mv myapps.xml ../myapps.xml
rm deezer-scope.labsin_deezer-scope.ini.in
rm deezer-scope.labsin_deezer-scope-settings.ini.in
rm myapps.xml.in
