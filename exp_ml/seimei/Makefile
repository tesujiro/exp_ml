all:

bayes:
	./random_name.sh 3000 >./train.txt
	./random_name.sh 1000 > ./eval.txt
	go run main_bayesian.go

clcnn_learn:
	./random_name.sh 10000 >./train.txt
	python clcnn_learn.py

clcnn_predict:
	./random_name.sh 1000 > ./eval.txt
	python clcnn_predict.py

