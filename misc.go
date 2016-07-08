package main

func coalesce(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// provided a number and a bounding range for that number, normalizes that
// number within that bounds (linearly). that is, it returns a float between 0
// and 1 representing i's position within the (min,max) range. Min and max are
// both inclusive.
func norm(i, min, max int) (n float64) {
	if i < min {
		return 0
	}
	if i > max {
		return 1
	}

	span := max - min
	reach := i - min

	return float64(reach) / float64(span)
}
