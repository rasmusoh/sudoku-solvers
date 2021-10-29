# Two sudoku solvers

This repo contains the same sudoku solver, written in python and go. Both uses the same basic approach: Constraint propagation using a matrix of sets, and when that fails, hypotheticals. The golang solution is approximately 240x faster (1M puzzles solved in 9.2 seconds, vs 10K puzzles in 22.0 seconds on my laptop) by using parallelism, efficient allocation patterns and by using bitflags instead of hash sets as the set implementation.
