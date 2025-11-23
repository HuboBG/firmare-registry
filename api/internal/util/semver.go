package util

import (
    "strconv"
    "strings"
)

// CompareSemver compares a and b as major.minor.patch.
// Returns 1 if a>b, -1 if a<b, 0 if equal.
func CompareSemver(a, b string) int {
    parse := func(v string) ([]int, bool) {
        parts := strings.Split(v, ".")
        if len(parts) < 1 || len(parts) > 3 {
            return nil, false
        }
        out := make([]int, 3)
        for i := 0; i < 3; i++ {
            if i < len(parts) {
                n, err := strconv.Atoi(parts[i])
                if err != nil {
                    return nil, false
                }
                out[i] = n
            }
        }
        return out, true
    }

    pa, oka := parse(a)
    pb, okb := parse(b)
    if !oka || !okb {
        if a == b {
            return 0
        }
        if a > b {
            return 1
        }
        return -1
    }

    for i := 0; i < 3; i++ {
        if pa[i] > pb[i] {
            return 1
        }
        if pa[i] < pb[i] {
            return -1
        }
    }
    return 0
}
