package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	. "github.com/jbrukh/bayesian"
)

const (
	TRAINING_DATA   = "./train.txt"
	EVALUATION_DATA = "./eval.txt"
)

const (
	CHAR_1 Class = "1"
	CHAR_2 Class = "2"
	CHAR_3 Class = "3"
	CHAR_4 Class = "4"
	CHAR_5 Class = "5"
)

func learn(c *Classifier) error {
	file, err := os.Open(TRAINING_DATA)
	if err != nil {
		return err
	}
	defer file.Close()
	r := csv.NewReader(file)
	r.Comma = ':'
	r.LazyQuotes = true

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		//for i, v := range record {
		//fmt.Printf("col[%v]:%v\t", i, v)
		//}
		DIVIDE_NUM := record[3]
		NAME := record[2]
		LEARN_DATA := []string{}
		for _, rn := range NAME {
			//fmt.Printf("[%v]:%c\t", i, rn)
			LEARN_DATA = append(LEARN_DATA, fmt.Sprintf("%c", rn))
		}
		//fmt.Printf("DIV=%v,LEARN_DATA=%v\n", DIVIDE_NUM, LEARN_DATA)
		c.Learn(LEARN_DATA, Class(DIVIDE_NUM))
	}

	return nil
}

func evaluate(c *Classifier) error {
	file, err := os.Open(EVALUATION_DATA)
	if err != nil {
		return err
	}
	defer file.Close()
	r := csv.NewReader(file)
	r.Comma = ':'
	r.LazyQuotes = true

	var right, wrong int
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		DIVIDE_NUM := record[3]
		NAME := record[2]
		LEARN_DATA := []string{}
		for _, rn := range NAME {
			//fmt.Printf("[%v]:%c\t", i, rn)
			LEARN_DATA = append(LEARN_DATA, fmt.Sprintf("%c", rn))
		}

		//scores, likely, _ := c.LogScores(LEARN_DATA)
		//fmt.Printf("name=%v\tscores=%v\tlikely=%v\n", NAME, scores, likely)

		correct, err := strconv.Atoi(DIVIDE_NUM)
		if err != nil {
			fmt.Printf("Number conversion error:%v", err)
			os.Exit(1)
		}
		probs, likely, _ := c.ProbScores(LEARN_DATA)
		if correct != likely+1 || correct != 2 {
			fmt.Printf("name=%v\tprobs=%8.8f\tlikely=%8.8d\n", NAME, probs, likely)
			wrong++
		} else {
			right++
		}
	}
	fmt.Printf(" Collect Ratio %v / %v = %v%\n", right, right+wrong, right*100/(right+wrong))
	return nil
}

func main() {
	classifier := NewClassifier(CHAR_1, CHAR_2, CHAR_3, CHAR_4, CHAR_5)
	if err := learn(classifier); err != nil {
		fmt.Println("Training ERROR :", err)
		os.Exit(1)
	}

	if err := evaluate(classifier); err != nil {
		fmt.Println("Evaluation ERROR :", err)
		os.Exit(1)
	}

	/*
		scores, likely, _ := classifier.LogScores([]string{"tall", "girl"})
		fmt.Printf("scores=%v\tlikely=%v\n", scores, likely)

		probs, likely, _ := classifier.ProbScores([]string{"tall", "girl"})
		fmt.Printf("probs=%8.8f\tlikely=%8.8d\n", probs, likely)
	*/
}
