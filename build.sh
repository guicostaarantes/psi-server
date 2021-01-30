mkdir -p /out/templates
for templ in  `find ./ -path "**/templates/**.html" -type f`
do
    cp $templ /out/templates/`basename $templ`
done
go build -o /out/main main.go