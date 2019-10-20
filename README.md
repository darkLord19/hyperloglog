# HyperLogLog
HyperLogLog is an algorithm for the count-distinct problem, approximating the number of distinct elements in a multiset.
Calculating the exact cardinality of a multiset requires an amount of memory proportional to the cardinality, which is impractical for very large data sets. 
Probabilistic cardinality estimators, such as the HyperLogLog algorithm, use significantly less memory than this, at the cost of obtaining only an approximation of the cardinality.