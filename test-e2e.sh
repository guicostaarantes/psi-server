mkdir -p ./e2e/templates
for templ in  `find ./ -path "**/templates/**.html" -type f`
do
    cp $templ ./e2e/templates/`basename $templ`
done
go test -coverpkg ./... -coverprofile ./e2e/cover.out ./e2e
go tool cover -html ./e2e/cover.out -o ./e2e/cover.html
rm -rf ./e2e/cover.out
rm -rf ./e2e/templates