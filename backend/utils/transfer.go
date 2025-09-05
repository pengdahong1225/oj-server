package utils

import (
	"oj-server/module/db"
	"oj-server/module/proto/pb"
)

func Transform(problem *db.Problem) *pb.Problem {
	return &pb.Problem{
		Id:          problem.ID,
		CreateAt:    problem.CreateAt.String(),
		Title:       problem.Title,
		Description: problem.Description,
		Level:       problem.Level,
		Tags:        SplitStringWithX(string(problem.Tags), "#"),
		CreateBy:    problem.CreateBy,
	}
}
func TransformList(list []db.Problem) []*pb.Problem {
	var problems []*pb.Problem
	for _, problem := range list {
		problems = append(problems, Transform(&problem))
	}
	return problems
}
