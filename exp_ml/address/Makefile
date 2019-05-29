all: prepare data
prepare:
	unzip -f ken_all.zip
	iconv -f CP932 -t UTF8 KEN_ALL.CSV | awk -F, '!/（|場合|、|〜/'> postcode.csv

learn:
	./random_address.sh 10000 >./train.txt
	python clcnn_learn_address.py

k2tf:
	-rm -r ./tensorflow_model
	python keras_to_tensorflow.py ./model.h5

predict:
	./random_address.sh 10000 > ./eval.txt
	go run clcnn_predict_address.go

