package service

import (
	"context"
	"k3s-nclink-apps/data-source/entity"
	"k3s-nclink-apps/model-manage-backend/mqtt"
	"log"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type adapterService struct {
	coll *mgm.Collection
}

var AdapterServ = &adapterService{
	coll: mgm.Coll(&entity.Adapter{}),
}

func (a *adapterService) Create(adapter *entity.Adapter) error {
	ModelServ.Lock()
	defer ModelServ.Unlock()
	model, err := ModelServ.FindByName(adapter.ModelName)
	if err != nil {
		return err
	}
	if err = a.coll.Create(adapter); err != nil {
		return err
	}
	model.Used++
	if err = ModelServ.update(model); err != nil {
		return err
	}
	ctx := context.Background()
	num, err := a.coll.EstimatedDocumentCount(ctx)
	if num <= 1 {
		_, err = a.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		})
	}
	return err
}

func (a *adapterService) Save(adapter *entity.Adapter, model *entity.Model) error {
	ModelServ.Lock()
	defer ModelServ.Unlock()
	if err := a.coll.Create(adapter); err != nil {
		return err
	}
	model.Used++
	if err := ModelServ.update(model); err != nil {
		return err
	}
	return a.coll.First(bson.M{"name": adapter.Name}, adapter)
}

// Find adapter
func (a *adapterService) FindById(id string) (*entity.Adapter, error) {
	ret := &entity.Adapter{}
	err := a.coll.FindByID(id, ret)
	return ret, err
}

func (a *adapterService) FindByName(name string) (*entity.Adapter, error) {
	ret := &entity.Adapter{}
	err := a.coll.First(bson.M{"name": name}, ret)
	return ret, err
}

func (a *adapterService) FindAll() ([]entity.Adapter, int64, error) {
	ret := []entity.Adapter{}
	err := a.coll.SimpleFind(&ret, bson.M{})
	if err != nil {
		return nil, 0, err
	}
	num, err := a.coll.EstimatedDocumentCount(context.Background())
	return ret, num, err
}

func (a *adapterService) FindByModelName(modelName string) ([]entity.Adapter, error) {
	ret := []entity.Adapter{}
	err := a.coll.SimpleFind(&ret, bson.M{"model_name": modelName})
	return ret, err
}

func (a *adapterService) delete(adapter *entity.Adapter) error {
	if err := a.coll.Delete(adapter); err != nil {
		return err
	}
	ModelServ.Lock()
	defer ModelServ.Unlock()
	model, err := ModelServ.FindByName(adapter.ModelName)
	if model.Used > 0 {
		model.Used--
		err = ModelServ.update(model)
	}
	return err
}

func (a *adapterService) DeleteByName(name string) error {
	adapter, err := a.FindByName(name)
	if err != nil {
		return err
	}
	return a.delete(adapter)
}

func (a *adapterService) DeleteById(id string) error {
	adapter, err := a.FindById(id)
	if err != nil {
		return err
	}
	return a.delete(adapter)
}

func (a *adapterService) update(adapter *entity.Adapter) error {
	return a.coll.Update(adapter)
}

func (a *adapterService) changeModel(adapter *entity.Adapter, modelName string) error {
	if adapter.ModelName == modelName {
		return nil
	}
	ModelServ.Lock()
	defer ModelServ.Unlock()
	newModel, err := ModelServ.FindByName(modelName)
	if err != nil {
		return err
	}
	model, _ := ModelServ.FindByName(adapter.ModelName)
	adapter.ModelName = modelName
	if err = a.update(adapter); err != nil {
		return err
	}
	if model.Used > 0 {
		model.Used--
		if err = ModelServ.update(model); err != nil {
			return err
		}
	}
	newModel.Used++
	return ModelServ.update(newModel)
}

func (a *adapterService) UpdateById(id string, in *entity.Adapter) (changed bool, err error) {
	adapter, err := a.FindById(id)
	if err != nil {
		return
	}
	if adapter.DevId == in.DevId && adapter.ModelName == in.ModelName {
		return
	}
	adapter.DevId = in.DevId
	if err = a.changeModel(adapter, in.ModelName); err != nil {
		return
	}
	mqtt.ResetModel(adapter.Name)
	return true, nil
}

func (a *adapterService) Rename(id, newName string) error {
	adapter, err := a.FindById(id)
	if err != nil {
		return err
	}
	oldName := adapter.Name
	if oldName != newName {
		adapter.Name = newName
		err = a.update(adapter)
	}
	return err
}

func (a *adapterService) ResetModel(adapters ...entity.Adapter) {
	for _, adapter := range adapters {
		log.Println("reset model on:", adapter.Name)
		mqtt.ResetModel(adapter.Name)
	}
}

func (a *adapterService) RenameModel(newName string, adapters ...entity.Adapter) error {
	for _, adapter := range adapters {
		adapter.ModelName = newName
		if err := a.coll.Update(&adapter); err != nil {
			return err
		}
	}
	return nil
}

func (a *adapterService) RenameModelFrom(oldName, newName string) error {
	adapters, err := a.FindByModelName(oldName)
	if err != nil {
		return err
	}
	return a.RenameModel(newName, adapters...)
}
