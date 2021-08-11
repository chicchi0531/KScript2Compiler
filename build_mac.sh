yacc_result=`goyacc -o compiler/parser.go compiler/parser.go.y`
if [ $yacc_result != ""]; then
    go run . testscript/$1.ks bin/$1.ksobj
    exit 0
fi
exit 1
