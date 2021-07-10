yacc_result=`goyacc -o compiler/parser.go compiler/parser.go.y`
if [ $yacc_result != ""]; then
    go run . $1
    exit 0
fi
exit 1
