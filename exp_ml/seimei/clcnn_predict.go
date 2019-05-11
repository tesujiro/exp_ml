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

	name_data, correct_data, err := load("./eval.txt")
	if err != nil {
		fmt.Printf("Error load data: %s\n", err)
		return
	}
	tensor, err := tf.NewTensor(name_data)
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
			model.Graph.Operation("dense_2/Sigmoid").Output(0),
		},
		nil,
	)

	if runErr != nil {
		fmt.Printf("TensorFlow runtime error: %s\n", runErr.Error())
		return
	}

	dense_out := output[0].Value().([][]float32)
	var right, wrong int
	for i, v := range dense_out {
		expected := int(correct_data[i])
		actual := argmax(v) + 1
		if expected == actual {
			right++
		} else {
			wrong++
			fmt.Printf("name=%v\tactual=%v\texpected=%v\n", code2name(name_data[i]), actual, expected)
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
func load(path string) ([][]float32, []float32, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	// input file format:
	// surname:givenname:fullname:divide numer
	// つのだ:ひろ:つのだひろ:3
	r := csv.NewReader(file)
	r.Comma = ':'
	r.LazyQuotes = true
	name_data := [][]float32{}
	correct_data := []float32{}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, nil, err
		}

		div := record[3]
		name := record[2]

		name_slice := []float32{}
		// name: string -> []float32
		for _, rn := range name {
			name_slice = append(name_slice, float32(rn))
		}
		for i := len(name_slice); i < 10; i++ {
			name_slice = append(name_slice, 0) // Padding
		}
		name_data = append(name_data, name_slice)

		correct, err := strconv.ParseFloat(div, 32)
		if err != nil {
			fmt.Printf("Number conversion error:%v", err)
			os.Exit(1)
		}
		correct_data = append(correct_data, float32(correct))
	}

	return name_data, correct_data, nil
}
