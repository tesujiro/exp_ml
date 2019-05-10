all:

bayes:
	./random_name.sh 3000 >./train.txt
	./random_name.sh 1000 > ./eval.txt
	go run main_bayesian.go

clcnn_learn:
	./random_name.sh 30000 >./train.txt
	python clcnn_learn.py

clcnn_predict:
	./random_name.sh 10000 > ./eval.txt
	python clcnn_predict.py

clcnn_test:
	./random_name.sh test > ./eval.txt
	python clcnn_predict.py

k2tf:
	python keras_to_tensorflow.py ./model.h5

