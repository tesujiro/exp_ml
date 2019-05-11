#!/usr/bin/env python
# coding: utf-8

"""
__doc__
General code to convert a trained keras model into an inference tensorflow model.
"""

import tensorflow as tf

model = tf.keras.models.load_model('./model.h5')
export_path = './tensorflow_model'

with tf.keras.backend.get_session() as sess:
    tf.saved_model.simple_save(
        sess,
        export_path,
        inputs={'input': model.input},
        outputs={t.name:t for t in model.outputs})

