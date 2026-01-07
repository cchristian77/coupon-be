package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PaginationTestSuite struct {
	suite.Suite
}

func (suite *PaginationTestSuite) Test_Limit() {
	type testCase struct {
		name           string
		expectedResult int
		inputLimit     int
	}

	testCases := []testCase{
		{
			name:           "get limit",
			expectedResult: 1,
			inputLimit:     1,
		},
		{
			name:           "default limit",
			expectedResult: DefaultLimit,
		},
	}

	for _, tc := range testCases {
		// Act
		var p Pagination
		p.SetLimit(tc.inputLimit)

		// Assert
		assert.Equal(suite.T(), tc.expectedResult, p.Limit())
	}
}

func (suite *PaginationTestSuite) Test_Page() {
	type testCase struct {
		name           string
		expectedResult int
		inputPage      int
	}

	testCases := []testCase{
		{
			name:           "get page",
			expectedResult: 1,
			inputPage:      1,
		},
		{
			name:           "default page",
			expectedResult: 1,
		},
	}

	for _, tc := range testCases {
		// Act
		var p Pagination
		p.SetPage(tc.inputPage)
		assert.Equal(suite.T(), tc.expectedResult, p.Page())
	}
}

func (suite *PaginationTestSuite) Test_Total() {
	type testCase struct {
		name           string
		expectedResult int64
		inputTotal     int64
	}

	testCases := []testCase{
		{
			name:           "get total",
			expectedResult: 10,
			inputTotal:     10,
		},
		{
			name:           "default total",
			expectedResult: 0,
		},
	}

	for _, tc := range testCases {
		// Act
		var p Pagination
		p.SetTotal(tc.inputTotal)

		// Assert
		assert.Equal(suite.T(), tc.expectedResult, p.Total())
	}
}

func (suite *PaginationTestSuite) Test_Offset() {
	type testCase struct {
		name           string
		expectedResult int
		input          struct {
			page  int
			limit int
		}
	}

	// offset = (page-1) * limit
	testCases := []testCase{
		{
			name:           "get offset page 1 limit 10",
			expectedResult: 0,
			input: struct {
				page  int
				limit int
			}{page: 1, limit: 10},
		},
		{
			name:           "get offset page 2 limit 5",
			expectedResult: 5,
			input: struct {
				page  int
				limit int
			}{page: 2, limit: 5},
		},
		{
			name:           "negative page considered as default page",
			expectedResult: 0,
			input: struct {
				page  int
				limit int
			}{page: -1, limit: 15},
		},
		{
			name:           "negative limit considered as default limit",
			expectedResult: 20,
			input: struct {
				page  int
				limit int
			}{page: 2, limit: -1},
		},
	}

	for _, tc := range testCases {
		// Act
		var p Pagination
		p.SetLimit(tc.input.limit)
		p.SetPage(tc.input.page)

		// Assert
		assert.Equal(suite.T(), tc.expectedResult, p.Offset())
	}
}

func (suite *PaginationTestSuite) Test_PageCount() {
	type testCase struct {
		name           string
		expectedResult int
		input          struct {
			limit int
			total int64
		}
	}

	// max page = round up(total / limit)
	testCases := []testCase{
		{
			name:           "limit 10 total 5",
			expectedResult: 1,
			input: struct {
				limit int
				total int64
			}{limit: 10, total: 5},
		},
		{
			name:           "limit 10 total 40",
			expectedResult: 4,
			input: struct {
				limit int
				total int64
			}{limit: 10, total: 40},
		},
		{
			name:           "limit 5 total 36",
			expectedResult: 8,
			input: struct {
				limit int
				total int64
			}{limit: 5, total: 36},
		},
	}

	for _, tc := range testCases {
		// Act
		var p Pagination
		p.SetLimit(tc.input.limit)
		p.SetTotal(tc.input.total)

		// Assert
		assert.Equal(suite.T(), tc.expectedResult, p.PageCount())
	}
}

func TestSuiteRunPagination(t *testing.T) {
	suite.Run(t, new(PaginationTestSuite))
}
