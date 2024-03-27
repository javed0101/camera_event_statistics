package constant

import (
	"net/http"

	"github.com/eencloud/goeen/api.v3"
)

var (
	ErrorNoEventFound = api.Reason{StatusCode: http.StatusBadRequest, Reason: "NoCameraEvent", Message: `No event found in redis`}

	ErrorInvalidParam      = api.Reason{StatusCode: http.StatusBadRequest, Reason: "InvalidQueryParams", Message: `cameraID or eventType is invalid`}
	ErrorQueryParamMissing = api.Reason{StatusCode: http.StatusBadRequest, Reason: "MissingQueryParam", Message: `cameraID or eventType is missing`}
	ErrorInvalidRequest    = api.Reason{StatusCode: http.StatusMethodNotAllowed, Reason: "InvalidRequest", Message: `Requsted resource could not found`}

	ErrorUnmarshallingJSON = api.Reason{StatusCode: http.StatusInternalServerError, Reason: "InternalServerError", Message: "Internal server error"}
	ErrorRedisDown         = api.Reason{StatusCode: http.StatusInternalServerError, Reason: "DatabaseDown", Message: "Unable to make connection with the DB"}
)
