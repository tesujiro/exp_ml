# -*- coding:utf-8 -*-
# 住所分割 with a Character-level CNN

import numpy as np
from keras.models import load_model, Model
from keras.layers import Input, Dense, Embedding, Reshape, Conv2D, MaxPooling2D, concatenate, BatchNormalization, Dropout
from keras.optimizers import Adam
from keras.callbacks import LearningRateScheduler
#import tensorflow as tf

MAX_LENGTH=30

def create_model(embed_size=64, max_length=MAX_LENGTH, filter_sizes=(2, 3, 4, 5, 6), filter_num=64):
    #inp = Input(shape=(max_length,),dtype=tf.int32)
    inp = Input(shape=(max_length,))
    emb = Embedding(0xffff, embed_size)(inp)
    emb_ex = Reshape((max_length, embed_size, 1))(emb)
    convs = []
    # Convolution2Dを複数通りかける
    for filter_size in filter_sizes:
        conv = Conv2D(filter_num, (filter_size, embed_size), activation="relu")(emb_ex)
        pool = MaxPooling2D(pool_size=(max_length - filter_size + 1, 1))(conv)
        convs.append(pool)
    print("loop finished")
    convs_merged = concatenate(convs)
    reshape = Reshape((filter_num * len(filter_sizes),))(convs_merged)
    fc1 = Dense(64, activation="relu")(reshape)
    bn1 = BatchNormalization()(fc1)
    do1 = Dropout(0.5)(bn1)
    out1 = Dense(MAX_LENGTH, activation='sigmoid', name='out1')(do1)
    out2 = Dense(MAX_LENGTH, activation='sigmoid', name='out2')(do1)
    out3 = Dense(MAX_LENGTH, activation='sigmoid', name='out3')(do1)
    model = Model(input=inp, output=out1)
    #model = Model(input=inp, outputs=[out1,out2,out3])
    return model

def load_data(filepath, max_length=MAX_LENGTH):
    l = []
    with open(filepath) as f:
        for line in f:
            address, pref, city, town = line.split(":")
            pref = int(pref) -1
            city = int(city) -1
            town = int(town) -1
            # divide name by characters
            # ord: character -> code point
            address = [ord(x) for x in address.strip()]
            address = address[:max_length]
            address_len = len(address)
            if address_len < max_length:
                # padding zeros
                address += ([0] * (max_length - address_len))
            l.append((pref, city, town, address))
    return l

def train(inputs, targets1, targets2, targets3, batch_size=100, epoch_count=100, max_length=MAX_LENGTH, model_filepath="./model.h5", learning_rate=0.001):

    # gradually decrease the learning rate
    start = learning_rate
    stop = learning_rate * 0.01
    learning_rates = np.linspace(start, stop, epoch_count)

    model = create_model(max_length=max_length)
    optimizer = Adam(lr=learning_rate)
    # categorical : more than 2 classes
    model.compile(loss='categorical_crossentropy',
                  optimizer=optimizer,
                  metrics=['accuracy'])

    model.fit(inputs,
              {'out1': targets1, 'out2': targets2, 'out3': targets3},
              epochs=epoch_count,
              batch_size=batch_size,
              verbose=1,
              validation_split=0.1,
              shuffle=True,
              callbacks=[
                  LearningRateScheduler(lambda epoch: learning_rates[epoch]),
              ])

    model.save(model_filepath)


if __name__ == "__main__":
    name_list = load_data('./train.txt')

    input_values = []
    target_values1 = []
    target_values2 = []
    target_values3 = []
    for target_value1, target_value2, target_value3, input_value in name_list:
        input_values.append(input_value)
        target_values1.append(target_value1)
        target_values2.append(target_value2)
        target_values3.append(target_value3)
    input_values = np.array(input_values)
    target_values1 = np.identity(MAX_LENGTH)[target_values1]
    target_values2 = np.identity(MAX_LENGTH)[target_values2]
    target_values3 = np.identity(MAX_LENGTH)[target_values3]
    print(input_values.shape)
    print(target_values1.shape)
    print(target_values2.shape)
    print(target_values3.shape)
    train(input_values, target_values1, target_value2, target_values3, epoch_count=50)
