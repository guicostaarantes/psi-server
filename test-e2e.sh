mkdir -p ./e2e/templates
for templ in  `find ./ -path "**/templates/**.html" -type f`
do
    cp $templ ./e2e/templates/`basename $templ`
done
go test ./e2e
rm -rf ./e2e/templates