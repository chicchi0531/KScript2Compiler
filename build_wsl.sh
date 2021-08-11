yacc_result=`goyacc.exe -o compiler/parser.go compiler/parser.go.y`
if [ $yacc_result != ""]; then
    go.exe build -o ks2compiler.exe app.go
    ./ks2compiler.exe testscript/$1.ks bin/$1.ksobj
    exit 0
fi
exit 1
