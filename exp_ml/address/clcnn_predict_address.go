package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

func main() {

	model, err := tf.LoadSavedModel("./tensorflow_model", []string{"serve"}, nil)

	if err != nil {
		fmt.Printf("Error loading saved model: %s\n", err.Error())
		return
	}

	defer model.Session.Close()

	address, correct_1, correct_2, correct_3, err := load("./eval.txt")
	if err != nil {
		fmt.Printf("Error load data: %s\n", err)
		return
	}
	tensor, err := tf.NewTensor(address)
	if err != nil {
		fmt.Printf("Error creating input tensor: %s\n", err.Error())
		return
	}

	// Print operation names
	/*
		for i, op := range model.Graph.Operations() {
			fmt.Printf("operation[%v]: %v\n", i, op.Name())
		}
	*/

	// You can get operation names from Keras model file.
	// See keras_tensor_flow.py
	output, runErr := model.Session.Run(
		map[tf.Output]*tf.Tensor{
			model.Graph.Operation("input_1").Output(0): tensor,
		},
		[]tf.Output{
			model.Graph.Operation("out1/Sigmoid").Output(0),
			model.Graph.Operation("out2/Sigmoid").Output(0),
			model.Graph.Operation("out3/Sigmoid").Output(0),
		},
		nil,
	)

	if runErr != nil {
		fmt.Printf("TensorFlow runtime error: %s\n", runErr.Error())
		return
	}

	output1 := output[0].Value().([][]float32)
	output2 := output[1].Value().([][]float32)
	output3 := output[2].Value().([][]float32)
	var right, wrong int
	for i, _ := range output1 {
		expected1 := int(correct_1[i])
		expected2 := int(correct_2[i])
		expected3 := int(correct_3[i])
		//actual := argmax(v) + 1
		actual1 := argmax(output1[i]) + 1
		actual2 := argmax(output2[i]) + 1
		actual3 := argmax(output3[i]) + 1
		if expected1 == actual1 && expected2 == actual2 && expected3 == actual3 {
			right++
		} else {
			wrong++
			if actual1 <= actual2 && actual2 <= actual3 {
				addr := []rune(code2name(address[i]))
				fmt.Printf("%v|%v|%v|%v\n", string(addr[0:actual1]), string(addr[actual1:actual2]), string(addr[actual2:actual3]), string(addr[actual3:]))
				//fmt.Printf("name=%v\tactual1=%v\tactual2=%v\tactual3=%v\n", code2name(address[i]), actual1, actual2, actual3)
			} else {
				fmt.Printf("name=%v\tactual1=%v\tactual2=%v\tactual3=%v\n", code2name(address[i]), actual1, actual2, actual3)
			}
		}
	}
	fmt.Printf("right: %v\twrong: %v\tratio: %v%%\n", right, wrong, float32(right*100)/float32(right+wrong))
}

func code2name(code []float32) string {
	var name string
	for _, ch := range code {
		name = fmt.Sprintf("%s%c", name, int(ch))
	}
	return name
}

func argmax(args []float32) int {
	index := 0
	max := args[index]
	for i, arg := range args {
		if arg > max {
			index = i
			max = arg
		}
	}
	return index
}

// load func returns name data and divide numbers.
func load(path string) ([][]float32, []float32, []float32, []float32, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer file.Close()

	// input file format:
	// surname:givenname:fullname:divide numer
	// 徳島県徳島市新浜町:3:6:9
	r := csv.NewReader(file)
	r.Comma = ':'
	r.LazyQuotes = true
	address_data := [][]float32{}
	correct_data1 := []float32{}
	correct_data2 := []float32{}
	correct_data3 := []float32{}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, nil, nil, nil, err
		}

		address := record[0]
		div1 := record[1]
		div2 := record[2]
		div3 := record[3]

		slice := []float32{}
		// address: string -> []float32
		for _, rn := range address {
			slice = append(slice, float32(rn))
		}
		for i := len(slice); i < 30; i++ {
			slice = append(slice, 0) // Padding
		}
		address_data = append(address_data, slice)

		correct1, err := strconv.ParseFloat(div1, 32)
		if err != nil {
			fmt.Printf("Number conversion error:%v", err)
			os.Exit(1)
		}
		correct_data1 = append(correct_data1, float32(correct1))

		correct2, err := strconv.ParseFloat(div2, 32)
		if err != nil {
			fmt.Printf("Number conversion error:%v", err)
			os.Exit(1)
		}
		correct_data2 = append(correct_data2, float32(correct2))

		correct3, err := strconv.ParseFloat(div3, 32)
		if err != nil {
			fmt.Printf("Number conversion error:%v", err)
			os.Exit(1)
		}
		correct_data3 = append(correct_data3, float32(correct3))
	}

	return address_data, correct_data1, correct_data2, correct_data3, nil
}
