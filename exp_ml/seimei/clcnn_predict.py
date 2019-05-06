# -*- coding:utf-8 -*-
# 姓名分割 with a Character-level CNN

import numpy as np
from keras.models import load_model, Model
from keras.layers import Input, Dense, Embedding, Reshape, Conv2D, MaxPooling2D, concatenate, BatchNormalization, Dropout

def load_data(filepath, max_length=10):
    name_list = []
    with open(filepath) as f:
        for l in f:
            surname, givenname, name, divide = l.split(":")
            divide = int(divide) -1
            # divide name by characters
            # ord: character -> code point
            name = [ord(x) for x in name.strip()]
            name = name[:max_length]
            name_len = len(name)
            if name_len < max_length:
                # padding zeros
                name += ([0] * (max_length - name_len))
            name_list.append((divide, name))
    return name_list

def predict(name_list, model_filepath="model.h5"):
    model = load_model(model_filepath)
    #ret = model.predict(name_list)
    ret = model.predict(name_list)
    return ret

if __name__ == "__main__":
    name_list = load_data('./eval.txt')

    input_values = []
    target_values = []
    for target_value, input_value in name_list:
        input_values.append(input_value)
        #target_values.append(target_value)
        #target_values.append(np.identity(5)[target_value])
        target_values = np.identity(5)[target_value]
        name = "".join([chr(i) for i in input_value])
        ret = predict(np.array([input_value]))
        expect = np.argmax(target_values)
        actual = np.argmax(ret)
        if (expect==actual):
            print("NAME:{}\tDIV={}".format(name,actual+1))
        else:
            print("NAME:{}\texpect:{}\tactual:{}".format(name,expect+1,actual+1))
        #print("".format(target_values))
        #print("          ret=={}".format(np.argmax(ret)))

    '''
    print("input_values={}".format(input_values))
    results = predict(input_values)
    for result in result:
        preint("result={}".format(result))

    right_count = 0
    wrong_count = 0
    '''

    
    '''
    raw_comment = "デートには最高やで！"
    comment = [ord(x) for x in raw_comment.strip().decode("utf-8")]
    comment = comment[:300]
    if len(comment) < 10:
        exit("too short!!")
    if len(comment) < 300:
        comment += ([0] * (300 - len(comment)))
    ret = predict(np.array([comment]))
    predict_result = ret[0][0]
    print "リア充度: {}%".format(predict_result * 100)
    '''

