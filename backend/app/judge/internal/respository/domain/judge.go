package domain

import (
	"github.com/sirupsen/logrus"
	"oj-server/module/model"
)

// UpdateUserSubmitRecord
// 1.更新用户提交记录表
// 2.更新(用户,题目)AC表
// 3.更新用户解题情况统计表
func (r *MysqlDB) UpdateUserSubmitRecord(record *model.SubmitRecord, level int32) error {
	tx := r.db_.Begin()
	if tx.Error != nil {
		logrus.Errorf("tx error: %v", tx.Error)
		return tx.Error
	}

	// todo 更新用户提交记录表
	/*
		insert into user_submit_record
		(uid, user_name, problem_id, problem_name, status, code, result, lang)
		values
		(?, ?, ?, ?, ?, ?, ?, ?);
	*/
	result := tx.Create(record)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		tx.Rollback()
		return result.Error
	}

	// todo 更新(用户,题目)AC表
	/*
		insert into user_solution(uid,problem_id)
		values(?,?)
	*/
	var repeatedAc = true
	if record.Status == "Accepted" {
		data := model.UserSolution{}
		result = tx.Where("uid=? and problem_id=?", record.Uid, record.ProblemID).Find(&data)
		if result.RowsAffected == 0 {
			repeatedAc = false
			data.Uid = record.Uid
			data.ProblemID = record.ProblemID

			result = tx.Create(&data)
			if result.Error != nil {
				logrus.Errorln(result.Error.Error())
				tx.Rollback()
				return result.Error
			}
		}
	}

	// todo 更新用户解题情况统计表
	data := model.Statistics{
		Uid: record.Uid,
	}
	result = tx.FirstOrCreate(&data)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		return result.Error
	}
	data.SubmitCount += 1
	if !repeatedAc && record.Status == "Accepted" {
		data.AccomplishCount += 1
		switch level {
		case 1:
			data.EasyProblemCount += 1
		case 2:
			data.MediumProblemCount += 1
		case 3:
			data.HardProblemCount += 1
		}
	}
	result = tx.Where("uid=?", data.Uid).Save(&data)
	if result.Error != nil {
		logrus.Errorln(result.Error.Error())
		tx.Rollback()
		return result.Error
	}

	tx.Commit()
	return nil
}
