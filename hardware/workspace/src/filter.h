/*
 * filter.h
 *
 *  Created on: 22 Mar 2021
 *      Author: bjs3118
 */

#ifndef FILTER_H_
#define FILTER_H_

#include "globals.h"

int float_to_fixed(float input);

float fixed_to_float(int input);

int fixed_mult(int x, int y);

void quantised_filt(int coef[], int buffer[], int x_read, int * quantised, int N);

void filt(float buffer[], int x_read, float * filtered, int N);

#endif /* FILTER_H_ */
