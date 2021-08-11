yacc_result=`goyacc -o compiler/parser.go compiler/parser.go.y`
if [ $yacc_result != ""]; then
    go build -o ks2compiler app.go
    ./ks2compiler testscript/$1.ks bin/$1.ksobj
    exit 0
fi
exit 1
