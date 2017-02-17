package server

import "net/http"

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

var allRoutes = routes{
	route{
		"GETExercise",
		"GET",
		"/v1/exercise/{id}",
		GETExercise,
	},
	route{
		"DeleteExercise",
		"DELETE",
		"/v1/exercise/{id}",
		DeleteExercise,
	},
	route{
		"PostExercise",
		"POST",
		"/v1/exercise",
		POSTExercise,
	},
	route{
		"GetExercises",
		"GET",
		"/v1/exercises",
		GETExercises,
	},
}
