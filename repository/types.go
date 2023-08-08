// This file contains types that are used in the repository layer.
package repository

type GetTestByIdInput struct {
	Id int
}

type GetTestByIdOutput struct {
	Name string
}
