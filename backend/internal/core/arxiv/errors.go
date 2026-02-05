package arxiv

import "errors"

var (
	// ErrFetchFailed indicates that fetching papers from arXiv failed.
	ErrFetchFailed = errors.New("failed to fetch papers from arXiv")

	// ErrSearchFailed indicates that searching papers failed.
	ErrSearchFailed = errors.New("failed to search papers")

	// ErrInvalidResponse indicates that the arXiv API returned an invalid response.
	ErrInvalidResponse = errors.New("invalid response from arXiv API")

	// ErrNotFound indicates that the requested paper was not found.
	ErrNotFound = errors.New("paper not found")
)

// IsFetchFailed checks if the error is ErrFetchFailed.
func IsFetchFailed(err error) bool { return errors.Is(err, ErrFetchFailed) }

// IsSearchFailed checks if the error is ErrSearchFailed.
func IsSearchFailed(err error) bool { return errors.Is(err, ErrSearchFailed) }

// IsInvalidResponse checks if the error is ErrInvalidResponse.
func IsInvalidResponse(err error) bool { return errors.Is(err, ErrInvalidResponse) }

// IsNotFound checks if the error is ErrNotFound.
func IsNotFound(err error) bool { return errors.Is(err, ErrNotFound) }
