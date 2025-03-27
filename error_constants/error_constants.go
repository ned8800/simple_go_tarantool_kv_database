package errorconstants

import "errors"

const (
	ErrLoadConfig       = "Error loading config"
	ErrStartServer      = "Error starting server"
	ErrShutdown         = "Error shutting down"
	ErrTarantoolConnect = "Error connecting to tarantool"
)

const (
	ErrInitializeConfig  = "Error initializing config"
	ErrUnmarshalConfig   = "Error unmarshalling config"
	ErrReadConfig        = "Error reading config"
	ErrReadEnvironment   = "Error reading .env file"
	ErrGetDirectory      = "Error getting directory"
	ErrDirectoryNotFound = "Error finding directory"
)

// handler's errors
const (
	ErrParseJSON          = "Error parsing JSON"
	ErrAlreadyExists      = "Already exists"
	ErrSendJSON           = "Error sending JSON"
	ErrSomethingWentWrong = "something went wrong"
	ErrBadPayload         = "bad payload"
)

var (
	ErrNotFoundById     = errors.New("no data found by that id")
	ErrKeyAlreadyExists = errors.New("this id already exists")
	ErrInsertValue      = errors.New("error while inserting value")
	ErrUpdateValue      = errors.New("error while updating value")
	ErrDeleteValue      = errors.New("error while deleting value")
)

const (
	ErrDuplicateKey = "Duplicate key exists"
)
