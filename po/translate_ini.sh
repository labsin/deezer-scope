cp ../deezer/deezer-scope.labsin_deezer-scope.ini.in ./
cp ../deezer/deezer-scope.labsin_deezer-scope-settings.ini.in ./
LC_ALL=C intltool-merge -d -u ./ deezer-scope.labsin_deezer-scope.ini.in deezer-scope.labsin_deezer-scope.ini
LC_ALL=C intltool-merge -d -u ./ deezer-scope.labsin_deezer-scope-settings.ini.in deezer-scope.labsin_deezer-scope-settings.ini
LC_ALL=C intltool-merge -x -u ./ myapps.xml.in myapps.xml
mv deezer-scope.labsin_deezer-scope.ini ../deezer/deezer-scope.labsin_deezer-scope.ini
mv deezer-scope.labsin_deezer-scope-settings.ini ../deezer/deezer-scope.labsin_deezer-scope-settings.ini
rm deezer-scope.labsin_deezer-scope.ini.in
rm deezer-scope.labsin_deezer-scope-settings.ini.in
sed -e 's/&#10;&#13; /\n/g' myapps.xml | sed -e 's/&#10;&#13;/\n/g' | sed -e 's/&lt;/\</g' | sed -e 's/&gt;/\>/g' > myapps.xml.tmp
mv myapps.xml.tmp myapps.xml
