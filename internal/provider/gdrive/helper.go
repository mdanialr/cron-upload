package gdrive

import "google.golang.org/api/googleapi"

// toField convert the given string type to googleapi.Field type.
func toField(field string) googleapi.Field {
	return googleapi.Field(field)
}
