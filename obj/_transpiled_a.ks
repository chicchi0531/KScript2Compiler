import "alice2_template"

func main() void{



ST(ih_alice,pos1,1)
__systemcall("ChangeName","アリス")
__systemcall("ShowMessage","立ち絵が表示されます<p>")

CLR(1)
__systemcall("ChangeName","")
__systemcall("PlayVoice","vo10001")
__systemcall("ShowMessage","立ち絵が消去されました<p>ああ<n>")
__systemcall("ShowMessage","ここは、<r=かんじ>漢字</r>がルビ<r=ひょうじ>表示</r>になります。<p>")
__systemcall("ShowMessage","変数テスト%s%sの%s<p>",hoge,fuga,hoge)



}
