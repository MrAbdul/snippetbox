package models

import "errors"

// ErrNoRecord you might be wondering why we’re returning the ErrNoRecord error from our SnippetModel.Get() method,
// instead of sql.ErrNoRows directly. The reason is to help encapsulate the model completely, so that our handlers aren’t
// concerned with the underlying datastore or reliant on datastore-specific errors (like sql.ErrNoRows) for its behavior.
var ErrNoRecord = errors.New("models: no matching record found")
