package ci

import (
	"github.com/revel/revel"
	o_ci "github.com/zze326/devops-helper/app/models/orm/ci"
	v_ci "github.com/zze326/devops-helper/app/models/view/ci"
	gormc "github.com/zze326/devops-helper/app/modules/gormc/app/controllers"
	"github.com/zze326/devops-helper/app/results"
	"gorm.io/gorm"
)

type CIEnvRef struct {
	gormc.Controller
}

// ListByPipelineID 根据流水线 ID 获取
func (c CIEnvRef) ListByPipelineID(pipelineID int) revel.Result {
	var envRefModels []*o_ci.EnvRef
	if err := c.DB.Preload("Env").Preload("Stages").Preload("Stages.Task").Preload("Stages.Tasks").Where("pipeline_id = ?", pipelineID).Order("sort").Find(&envRefModels).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(envRefModels)
}

// Save 保存
func (c CIEnvRef) Save(req v_ci.SaveEnvRefsReq) revel.Result {
	if err := c.DB.Transaction(func(tx *gorm.DB) error {
		for _, envRefItem := range req.EnvRefs {
			// 处理 env
			envRefModel := new(o_ci.EnvRef)
			if envRefItem.ID == 0 {
				envRefModel.EnvID = envRefItem.EnvID
				envRefModel.PipelineID = req.PipelineID
				envRefModel.Sort = envRefItem.Sort
			} else {
				if err := tx.First(envRefModel, envRefItem.ID).Error; err != nil {
					return err
				}
				envRefModel.Sort = envRefItem.Sort
			}
			if err := tx.Save(envRefModel).Error; err != nil {
				return err
			}

			// 处理 stages
			for _, stageItem := range envRefItem.Stages {
				stageModel := new(o_ci.Stage)
				if stageItem.ID > 0 {
					if err := tx.First(stageModel, stageItem.ID).Error; err != nil {
						return err
					}
				}
				stageModel.Name = stageItem.Name
				stageModel.Parallel = stageItem.Parallel
				stageModel.Sort = stageItem.Sort
				stageModel.EnvRefID = envRefModel.ID

				if err := tx.Save(stageModel).Error; err != nil {
					return err
				}

				if stageItem.Parallel {
					for _, taskItem := range stageItem.Tasks {
						taskModel := new(o_ci.Task)
						if taskItem.ID > 0 {
							if err := tx.First(taskModel, taskItem.ID).Error; err != nil {
								return err
							}
						}
						taskModel.Type = taskItem.Type
						taskModel.Content = taskItem.Content
						taskModel.Url = taskItem.Url
						taskModel.Branch = taskItem.Branch
						taskModel.StageID = stageModel.ID
						if err := tx.Save(taskModel).Error; err != nil {
							return err
						}
					}
				} else {
					taskModel := new(o_ci.Task)
					if stageItem.Task.ID > 0 {
						if err := tx.First(taskModel, stageItem.Task.ID).Error; err != nil {
							return err
						}
					}
					taskModel.Type = stageItem.Task.Type
					taskModel.Content = stageItem.Task.Content
					taskModel.Url = stageItem.Task.Url
					taskModel.Branch = stageItem.Task.Branch
					if err := tx.Save(taskModel).Error; err != nil {
						return err
					}

					stageModel.TaskID = taskModel.ID
					if err := tx.Updates(stageModel).Error; err != nil {
						return err
					}
				}

			}

		}
		return nil
	}); err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}
