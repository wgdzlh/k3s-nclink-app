package conv

import (
	"k3s-nclink-apps/config-distribute/models/entity"
	pb "k3s-nclink-apps/configmodel"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func MessageToString(msg proto.Message) (string, error) {
	json, err := protojson.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(json), err
}

func DbModelToWireModel(in *entity.Model) (*pb.Model, error) {
	ret := &pb.Model{}
	err := protojson.Unmarshal([]byte(in.Def), ret)
	if err != nil {
		return nil, err
	}
	ret.Name = in.Name
	return ret, nil
}