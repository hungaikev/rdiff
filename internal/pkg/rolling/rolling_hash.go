package rolling

import (
	"context"

	"go.opentelemetry.io/otel"
)

/**

The Hash function calculates the rolling hash value for a given slice of bytes using the Rolling Hash algorithm. Here is a step-by-step breakdown of how the function works:

1. It calculates the length of the input data, n.
2. It initializes the rolling hash value, hash, to 0.
3. It creates an array of powers of two, powers, with a length of n. The first element in the array is set to 1.
4. It pre-calculates the powers of two for each element in the powers array. The ith element is set to the value of the (i-1)th element multiplied by 2.
5. It iterates over each byte in the input data. For each iteration, it adds the current byte to the rolling hash value by multiplying it by the corresponding power of two in the powers array and adding it to hash.
6. It returns the final rolling hash value, hash.

*/

var tracer = otel.Tracer("rolling")

// Hash calculates the rolling hash value for the given data using the Rolling Hash algorithm
func Hash(ctx context.Context, data []byte) uint64 {
	ctx, span := tracer.Start(ctx, "rolling.Hash")
	defer span.End()

	// the length of the data
	n := len(data)

	// the initial value for the rolling hash
	hash := uint64(0)

	// an array of powers of two
	powers := make([]uint64, n)
	powers[0] = 1

	// pre-calculate the powers of two
	for i := 1; i < n; i++ {
		powers[i] = powers[i-1] * 2
	}

	// calculate the rolling hash value
	for i := 0; i < n; i++ {
		// add the current byte to the rolling hash
		hash += uint64(data[i]) * powers[i]
	}

	return hash
}

// FromRollingHash converts the given rolling hash value back into the original data
func FromRollingHash(ctx context.Context, hash uint64, n int) []byte {
	ctx, span := tracer.Start(ctx, "rolling.FromRollingHash")
	defer span.End()

	// the original data
	data := make([]byte, n)

	// an array of powers of two
	powers := make([]uint64, n)
	powers[0] = 1

	// pre-calculate the powers of two
	for i := 1; i < n; i++ {
		powers[i] = powers[i-1] * 2
	}

	// retrieve the original data
	for i := n - 1; i >= 0; i-- {
		// divide the hash value by the power of two
		// to get the current byte
		data[i] = byte(hash / powers[i])

		// update the hash value by subtracting the
		// byte value multiplied by the power of two
		hash -= uint64(data[i]) * powers[i]
	}

	return data
}
