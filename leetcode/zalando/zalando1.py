""""

You are given an implementation of a function: def solution(A, X) 
This function, when given an array A of N integers, sorted in non-decreasing 
order, and some integer X, looks for X in A. If X is present in A, then the
 function returns position of (some) occurrence of X in A. Otherwise, the
  function returns −1. For example, given the following array: A[0] = 
  1 A[1] = 2 A[2] = 5 A[3] = 9 A[4] = 9 and X = 5, the function should 
  return 2, as A[2] = 5. The attached code is still incorrect for some inputs.
   Despite the error(s), the code may produce a correct answer for the example 
   test cases. The goal of the exercise is to find and fix the bug(s) in 
   the implementation. You can modify at most three lines. Assume that: 
   N is an integer within the range [0..100,000]; each element of array
    A is an integer within the range [−2,000,000,000..2,000,000,000]; 
    array A is sorted in non-decreasing order; X is an integer within 
    the range [−2,000,000,000..2,000,000,000]. 
"""

def solution(A, X):
    N = len(A)
    if N == 0:
        return -1
    l = 0
    r = N - 1
    while l <= r:
        m = (l + r) // 2
        if A[m] > X:
            r = m - 1
        else:
            l = m + 1
    if l > 0 and A[l-1] == X:
        return l-1
    return -1

def solution(A, X):
    N = len(A)
    if N == 0:
        return -1
    l = 0
    r = N - 1
    while l <= r:  # Linha alterada
        m = (l + r) // 2
        if A[m] > X:
            r = m - 1
        else:
            l = m + 1  # Linha alterada
    if r >= 0 and A[r] == X:  # Linha alterada
        return r
    return -1

A = [1, 2, 5, 9, 9]
X = 5

print(solution(A, X))


# def solution(A, X):
#     N = len(A)
#     if N == 0:
#         return -1
#     l = 0
#     r = N - 1
#     while l <= r:
#         m = (l + r) // 2
#         if A[m] > X:
#             r = m - 1
#         else:
#             l = m + 1  
#     if A[l - 1] == X:
#         return l - 1
#     return -1

#ORIGINAL
# def solution(A, X):
#     N = len(A)
#     if N == 0:
#         return -1
#     l = 0
#     r = N - 1
#     while l < r:
#         m = (l + r) // 2
#         if A[m] > X:
#             r = m - 1
#         else:
#             l = m
#     if A[l] == X:
#         return l
#     return -1
