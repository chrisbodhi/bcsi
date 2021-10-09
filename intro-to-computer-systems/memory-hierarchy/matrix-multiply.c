/*
Naive code for multiplying two matrices together.

There must be a better way!
*/

#include <stdio.h>
#include <stdlib.h>

/*
  A naive implementation of matrix multiplication.

  DO NOT MODIFY THIS FUNCTION, the tests assume it works correctly, which it
  currently does
*/
void matrix_multiply(double **C, double **A, double **B, int a_rows, int a_cols,
                     int b_cols) {
  for (int i = 0; i < a_rows; i++) {
    for (int j = 0; j < b_cols; j++) {
      C[i][j] = 0;
      for (int k = 0; k < a_cols; k++) {
        C[i][j] += A[i][k] * B[k][j];
      }
    }
  }
}

void fast_matrix_multiply(double **c, double **a, double **b, int a_rows) {
  int a_cols = a_rows;
  int b_cols = a_rows;
  for (int i = 0; i < a_rows; i++) {
    double* ai = a[i];
    for (int j = 0; j < b_cols; j++) {
      c[i][j] = 0;
      for (int k = 0; k < a_cols; k++) {
        /* printf("i: %d, j: %d, k: %d\n", i, j, k); */
        c[i][j] += ai[k] * b[k][j];
      }
    }
  }
}

/* Local Variables: */
/* rmsbolt-command: "clang -O1" */
/* End: */
