#!/bin/bash

INFILE=./postcode.csv
REPEAT=${1:-1000}

# 郵便番号データダウンロード
# https://www.post.japanpost.jp/zipcode/download.html

# 国土数値情報
# http://nlftp.mlit.go.jp/index.html
# 
# 位置参照情報ダウンロードサービス
# http://nlftp.mlit.go.jp/cgi-bin/isj/dls/_choose_method.cgi
#

conv()
{
    iconv -f CP932 -t UTF8 $1 | sed 's/'$'\r//g'
}

random()
{
    gawk -F, -v REPEAT=$1 '
    BEGIN{
	srand()
    }
    {
	address[NR]=$0
    }
    END{
    	for(i=1;i<=REPEAT;i++){
	    j=int(rand()*length(address))
	    print address[j]
	}
    }'
}

extract_col()
{
    gawk -F, '
    BEGIN{
    	OFS=":"
    }
    {
	print $7 $8 $9,length($7),length($7 $8),length($7 $8 $9)
    }'
}

if [ "$1" == "test" ]; then
    cat <<EOF
EOF
    exit 0
fi

cat $INFILE | sed -e "s/\"//g" | random $REPEAT | extract_col
