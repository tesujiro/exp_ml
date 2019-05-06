# -*- coding:utf-8 -*-
# 姓名分割 with a Character-level CNN

import numpy as np
from keras.models import load_model, Model
from keras.layers import Input, Dense, Embedding, Reshape, Conv2D, MaxPooling2D, concatenate, BatchNormalization, Dropout
from keras.optimizers import Adam
from keras.callbacks import LearningRateScheduler

def create_model(embed_size=256, max_length=10, filter_sizes=(2, 3, 4, 5), filter_num=64):
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
    fc2 = Dense(5, activation='sigmoid')(do1)
    model = Model(input=inp, output=fc2)
    return model

def load_data(filepath, max_length=10):
    name_list = []
    with open(filepath) as f:
        for l in f:
            surname, givenname, name, divide = l.split(":")
            divide = int(divide) -1
            # divide name by characters
            # ord: character -> code point
            #name = [ord(x) for x in name.strip().decode("utf-8")]
            name = [ord(x) for x in name.strip()]
            name = name[:max_length]
            name_len = len(name)
            if name_len < max_length:
                # padding zeros
                name += ([0] * (max_length - name_len))
            name_list.append((divide, name))
    return name_list

def train(inputs, targets, batch_size=100, epoch_count=100, max_length=10, model_filepath="./model.h5", learning_rate=0.001):

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

    model.fit(inputs, targets,
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
    target_values = []
    for target_value, input_value in name_list:
        input_values.append(input_value)
        target_values.append(target_value)
    input_values = np.array(input_values)
    print(input_values)
    target_values = np.identity(5)[target_values]
    train(input_values, target_values, epoch_count=50)
