#!/bin/bash

INFILE=./JINMEI30.TXT
#OUTFILE=./name_dictionary.txt
REPEAT=${1:-1000}

conv()
{
    iconv -f CP932 -t UTF8 $1 | sed 's/'$'\r//g'
}

format()
{
    gawk '
    BEGIN{
    	OFS=":"
    }
    NF==2 {
	KANA=$1
	n=split($2,a,":")
	if(n!=2) next;
	KANJI=a[1]
	gsub(/"/,"",KANJI)
	KIND=a[2]
	if(KANJI==""){
	    KANJI="ERROR NR:" NR
	}
	if(KANA==""){
	    KANA="ERROR NR:" NR
	}
	if(KIND=="姓"||KIND=="名"){
	    print KIND,KANJI,KANA
	}
    }'
}

random_name()
{
    gawk -F: -v REPEAT=$1 '
    {
	KIND=$1
	NAME=$2
	KANA=$3
	if(KIND=="姓"){
	    SNAME_COUNT++
	    SNAME_KANJI[SNAME_COUNT]=NAME
	    SNAME_KANA[SNAME_COUNT]=KANA
	} else {
	    GNAME_COUNT++
	    GNAME_KANJI[GNAME_COUNT]=NAME
	    FNAME_KANA[GNAME_COUNT]=KANA
	}
    }
    END{
    	OFS=":"
	srand()
	for(i=1;i<=REPEAT;i++){
	    RND_SNAME=SNAME_KANJI[int(rand()*SNAME_COUNT)+1]
	    RND_GNAME=GNAME_KANJI[int(rand()*GNAME_COUNT)+1]
	    FULLNAME=RND_SNAME RND_GNAME
	    DIV=length(RND_SNAME)
	    print RND_SNAME,RND_GNAME,FULLNAME,DIV
	}

    }'
}

if [ "$1" == "test" ]; then
    cat <<EOF
::つのだ☆ひろ:3
::子門前太郎:3
::小宮山和太郎:3
::斉藤禿頭:2
::齋藤禿頭:2
EOF
    exit 0
fi

unzip -qq ${INFILE}.zip 
conv $INFILE | format | random_name $REPEAT
rm $INFILE
