/*
Package types contains global definitions that are shared between the agent and API.
Has no dependencies and is easily importable for any Golang frontend.
*/
package types

// generic type aliases
// these just make code more readable
// lets you see which specific string is being passed where

// TeamNum is the identifier of a team.
// They start from 1 and go up to the configure number of teams
// [1, num_teams].
type TeamNum uint32

// UserListName is the unique identifier of a userlist.
type UserListName string

// ServiceName is the identifier of a service.
// This is unique under a single box, but different boxes can have services with the same ID.
type ServiceName string

// BoxName is the unique identifier of a box.
type BoxName string

// CheckName is the identifier of a check.
// This has no uniqueness constraint.
// It has no meaning to the API; is only used by the agent.
type CheckName string
