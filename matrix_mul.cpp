
#include <cstdlib> // For rand() and srand()
#include <ctime>   // For time()
#include <iomanip> // For std::setw and std::setfill
#include <iostream>
#include <vector>

const int SIZE = 30; // Changeable matrix size

// Function to print a matrix with improved formatting
void printMatrix(const std::vector<std::vector<int>> &matrix) {
  for (const auto &row : matrix) {
    std::cout << "| ";
    for (const auto &elem : row) {
      std::cout << std::setw(4) << elem << " "; // Align columns with setw
    }
    std::cout << "|\n";
  }
  std::cout << std::endl;
}

int main() {
  // Seed random number generator
  srand(static_cast<unsigned>(time(0)));

  // Initialize matrices A, B, and C
  std::vector<std::vector<int>> A(SIZE, std::vector<int>(SIZE));
  std::vector<std::vector<int>> B(SIZE, std::vector<int>(SIZE));
  std::vector<std::vector<int>> C(SIZE, std::vector<int>(SIZE, 0));

  // Fill matrices A and B with random numbers
  for (int i = 0; i < SIZE; ++i) {
    for (int j = 0; j < SIZE; ++j) {
      A[i][j] = rand() % 10; // Random numbers from 0 to 9
      B[i][j] = rand() % 10;
    }
  }

  // Display Matrix A
  std::cout << "Matrix A:\n";
  printMatrix(A);

  // Display Matrix B
  std::cout << "Matrix B:\n";
  printMatrix(B);

  // Multiply matrices A and B, store result in C
  for (int i = 0; i < SIZE; ++i) {
    for (int j = 0; j < SIZE; ++j) {
      for (int k = 0; k < SIZE; ++k) {
        C[i][j] += A[i][k] * B[k][j];
      }
    }
  }

  // Display Matrix C
  std::cout << "Matrix C (A * B):\n";
  printMatrix(C);

  return 0;
}
