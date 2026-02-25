//package tasks




//func DeleteTask(ctx context.Context, id uint) error {
//	_, err := gorm.G[models.Task](postgres.Db).Where("id = ?", id).Delete(ctx)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func UpdateTask(ctx context.Context, id uint, task *TaskData) error {
//	updates := models.Task{
//		Name:        task.Name,
//		Description: task.Description,
//		Status:      task.Status,
//		Priority:    task.Priority,
//	}
//
//	_, err := gorm.G[models.Task](postgres.Db).Where("id = ?", id).Updates(ctx, updates)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func GetTask(ctx context.Context, id uint) (*TaskData, error) {
//	task, err := gorm.G[models.Task](postgres.Db).Where("id =?", id).First(ctx)
//	if err != nil {
//		return nil, err
//	}
//	return &TaskData{
//		Name:        task.Name,
//		Description: task.Description,
//		Status:      task.Status,
//		Priority:    task.Priority,
//	}, nil
//}
