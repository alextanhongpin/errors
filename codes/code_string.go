// Code generated by "stringer -type=Code -linecomment"; DO NOT EDIT.

package codes

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[unknown-0]
	_ = x[Aborted-1]
	_ = x[BadRequest-2]
	_ = x[Canceled-3]
	_ = x[Conflict-4]
	_ = x[DataLoss-5]
	_ = x[DeadlineExceeded-6]
	_ = x[Exists-7]
	_ = x[Forbidden-8]
	_ = x[Internal-9]
	_ = x[NotFound-10]
	_ = x[NotImplemented-11]
	_ = x[OutOfRange-12]
	_ = x[PreconditionFailed-13]
	_ = x[TooManyRequests-14]
	_ = x[Unauthorized-15]
	_ = x[Unavailable-16]
	_ = x[Unknown-17]
}

const _Code_name = "unknownabortedbad_requestcanceledconflictdata_lossdeadline_exceededexistsforbiddeninternalnot_foundnot_implementedout_of_rangeprecondition_failedtoo_many_requestsunauthorizedunavailableunknown"

var _Code_index = [...]uint8{0, 7, 14, 25, 33, 41, 50, 67, 73, 82, 90, 99, 114, 126, 145, 162, 174, 185, 192}

func (i Code) String() string {
	if i < 0 || i >= Code(len(_Code_index)-1) {
		return "Code(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Code_name[_Code_index[i]:_Code_index[i+1]]
}