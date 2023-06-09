// Code generated by "stringer -type=Provider"; DO NOT EDIT.

package cloudprovider

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Unknown-0]
	_ = x[AWS-1]
	_ = x[Azure-2]
	_ = x[GCP-3]
	_ = x[OpenStack-4]
	_ = x[QEMU-5]
}

const _Provider_name = "UnknownAWSAzureGCPOpenStackQEMU"

var _Provider_index = [...]uint8{0, 7, 10, 15, 18, 27, 31}

func (i Provider) String() string {
	if i >= Provider(len(_Provider_index)-1) {
		return "Provider(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Provider_name[_Provider_index[i]:_Provider_index[i+1]]
}
